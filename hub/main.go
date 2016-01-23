package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/boltdb/bolt"
)

const hubUsage = `
    JET/Hub v0.4 (http://jeelabs.org/2016/01/overcoming-jet-lag/)

    Usage: /path/to/hub ?options...?
`

type Event struct {
	topic   string
	payload []byte
}

func main() {
	adminFlag := flag.String("admin", "", "connect as admin to a running hub")
	dataStore := flag.String("data", "store.db", "data store file name & path")
	mqttPort := flag.String("mqtt", "tcp://localhost:1883", "MQTT server port")
	httpPort := flag.String("http", "localhost:8947", "HTTP server port")
	flag.Parse()

	// check for special admin mode, used by the "jet" wrapper script
	if *adminFlag != "" {
		adminCmd(connectToHub("admin", *adminFlag, false))
		return
	}

	// due to the above, "--help" isn't very user-friendly, use "help" instead
	if flag.Arg(0) == "help" {
		fmt.Println(hubUsage)
		return
	}

	// normal hub startup begins here, with a log entry
	log.Print(append([]string{"JET/Hub"}, os.Args[1:]...))

	quit := make(chan struct{})

	// connect to MQTT and wait for it before doing anything else
	hub := connectToHub("hub", *mqttPort, true)
	defer hub.Disconnect(250)

	// send one message every second, on the second
	go sendHeartbeat(hub, "hub/1hz")

	// open the persistent data store
	log.Println("opening data store:", *dataStore)
	options := bolt.Options{Timeout: time.Second}
	db, err := bolt.Open(*dataStore, 0600, &options)
	if err != nil {
		log.Fatalln("db:", err)
	}
	defer db.Close()

	// look for serial device(s) and listen to them
	devChanges := topicAsEvents(hub, "serial/+")
	go listenToDevices(devChanges)

	// the default is to start up the built-in HTTP server
	if *httpPort != "" {
		go func() {
			defer close(quit)
			startHttpServer(*httpPort)
		}()
	}

	<-quit // hang around until something serious happens
}

func connectToHub(clientName, port string, retain bool) *mqtt.Client {
	// add a "fairly random" 6-digit suffix to make the client name unique
	nanos := time.Now().UnixNano()
	clientId := fmt.Sprintf("%s/%06d", clientName, nanos % 1e6)

	options := mqtt.NewClientOptions()
	options.AddBroker(port)
	options.SetClientID(clientId)
	options.SetBinaryWill("jet/" + clientId, nil, 1, retain)
	client := mqtt.NewClient(options)

	if t := client.Connect(); t.Wait() && t.Error() != nil {
		log.Fatal(t.Error())
	}

	log.Println("connected as", clientId, "to", port)

	// register as jet client, cleared on disconnect by the will
	t := client.Publish("jet/" + clientId, 1, retain, "{}")
	if t.Wait() && t.Error() != nil {
		log.Fatal(t.Error())
	}

	return client
}

func topicAsEvents(hub *mqtt.Client, pattern string) chan Event {
	feed := make(chan Event)

	t := hub.Subscribe(pattern, 0, func(hub *mqtt.Client, msg mqtt.Message) {
		feed <- Event{
			topic:   string(msg.Topic()),
			payload: msg.Payload(),
		}
	})
	if t.Wait() && t.Error() != nil {
		log.Fatal(t.Error())
	}

	return feed
}

func startHttpServer(port string) {
	http.HandleFunc("/bar",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", r.URL.Path)
		})

	certFile := os.Getenv("HUB_HTTP_CERT")
	keyFile := os.Getenv("HUB_HTTP_KEY")

	if certFile != "" && keyFile != "" {
		log.Println("starting HTTPS (TLS) server at", port)
		log.Fatal(http.ListenAndServeTLS(port, certFile, keyFile, nil))
	} else {
		log.Println("starting HTTP server at", port)
		log.Fatal(http.ListenAndServe(port, nil))
	}
}

func sendHeartbeat(hub *mqtt.Client, topic string) {
	for {
		time.Sleep(time.Duration(1e9 - time.Now().UnixNano()%1e9))

		// publish the heartbeat msg if it's within 25ms of the second mark
		millis := time.Now().UnixNano() / 1e6
		if millis%1000 < 25 {
			payload := fmt.Sprintf("%d", millis)
			t := hub.Publish(topic, 0, false, payload)
			if t.Wait() && t.Error() != nil {
				log.Print(t.Error())
			}
		} else {
			log.Println("missed heartbeat:", millis)
		}
		log.Println("heartbeat:", millis)
	}
}
