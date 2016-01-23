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

const VERSION = "0.4.0"

var hubUsage = fmt.Sprintf(`
    JET/Hub v%s (http://jeelabs.org/2016/01/overcoming-jet-lag/)

    Usage: /path/to/hub ?options...?
`, VERSION)

var hub *mqtt.Client

type Event struct {
	Topic   string
	Payload []byte
}

func main() {
	adminFlag := flag.String("admin", "", "connect as admin to a running hub")
	dataStore := flag.String("data", "store.db", "data store file name & path")
	mqttPort := flag.String("mqtt", "tcp://localhost:1883", "MQTT server port")
	httpPort := flag.String("http", "localhost:8947", "HTTP server port")
	flag.Parse()

	// check for special admin mode, used by the "jet" wrapper script
	if *adminFlag != "" {
		hub = connectToHub("admin", *adminFlag, false)
		adminCmd()
		return
	}

	// due to the above, "--help" isn't very user-friendly, use "help" instead
	if flag.Arg(0) == "help" {
		fmt.Println(hubUsage)
		return
	}

	// normal hub startup begins here, with a log entry
	log.Print(append([]string{"JET/Hub v" + VERSION}, os.Args[1:]...))

	quit := make(chan struct{})

	// connect to MQTT and wait for it before doing anything else
	hub = connectToHub("hub", *mqttPort, true)
	defer hub.Disconnect(250)

	// open the persistent data store
	log.Println("opening data store:", *dataStore)
	options := bolt.Options{Timeout: time.Second}
	db, err := bolt.Open(*dataStore, 0600, &options)
	if err != nil {
		log.Fatalln("db:", err)
	}
	defer db.Close()

	// send one message every second, on the second
	go sendHeartbeat("hub/1hz")

	// copy each incoming "logger/<x>" message to "logger/<x>/<millis>"
	go timestampRepeater(topicsAsEvents("logger/+"))

	// listen to serial device requests
	go processSerialRequests(topicsAsEvents("serial/+"))

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
	clientId := fmt.Sprintf("%s/%06d", clientName, nanos%1e6)

	options := mqtt.NewClientOptions()
	options.AddBroker(port)
	options.SetClientID(clientId)
	options.SetBinaryWill("jet/"+clientId, nil, 1, retain)
	client := mqtt.NewClient(options)

	if t := client.Connect(); t.Wait() && t.Error() != nil {
		log.Fatal(t.Error())
	}

	if retain {
		log.Println("connected as", clientId, "to", port)
	}

	// register as jet client, cleared on disconnect by the will
	t := client.Publish("jet/"+clientId, 1, retain, "{}")
	if t.Wait() && t.Error() != nil {
		log.Fatal(t.Error())
	}

	return client
}

func topicsAsEvents(pattern string) chan Event {
	feed := make(chan Event)

	t := hub.Subscribe(pattern, 0, func(hub *mqtt.Client, msg mqtt.Message) {
		feed <- Event{
			Topic:   string(msg.Topic()),
			Payload: msg.Payload(),
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

func sendHeartbeat(topic string) {
	for {
		time.Sleep(time.Duration(1e9 - time.Now().UnixNano()%1e9))

		// publish the heartbeat msg if it's within 25ms of the second mark
		millis := time.Now().UnixNano() / 1e6
		if millis%1000 < 25 {
			publish(topic, []byte(fmt.Sprintf("%d", millis)), false)
		} else {
			log.Println("missed heartbeat:", millis)
		}
	}
}

func publish(topic string, payload []byte, retain bool) {
	t := hub.Publish(topic, 0, retain, payload)
	if t.Wait() && t.Error() != nil {
		log.Print(t.Error())
	}
}
