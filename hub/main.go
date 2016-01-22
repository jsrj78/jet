package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/surge/glog"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

const hubUsage = `
    JET/Hub v0.4 (http://jeelabs.org/2016/01/overcoming-jet-lag/)

    Usage: /path/to/hub -logtostderr
`

var (
	adminFlag = flag.String("admin", "", "connect as admin to a running hub")
	dataStore = flag.String("data", "storage.db", "data store file name & path")
	mqttPort  = flag.String("mqtt", "localhost:1883", "MQTT server port")
	extServer = flag.Bool("external", false, "use an external MQTT server")
	httpPort  = flag.String("http", "localhost:8947", "HTTP server port")
)

type Event struct {
	topic   string
	payload []byte
}

func main() {
	flag.Parse()
	defer glog.Flush()

	// check for special admin mode, used by the "jet" wrapper script
	if *adminFlag != "" {
		adminCmd(connectToHub("admin", *adminFlag))
		return
	}

	// due to the above, "--help" isn't very user-friendly, use "help" instead
	if flag.Arg(0) == "help" {
		fmt.Println(hubUsage)
		return
	}

	// normal hub startup begins here, with a log entry
	glog.Info(append([]string{"JET/Hub"}, os.Args[1:]...))

	quit := make(chan struct{})

	// the default is to start the built-in MQTT server
	if !*extServer {
		go func() {
			defer close(quit)
			startMqttServer(*mqttPort)
		}()
	}

	// connect to MQTT and wait for it before doing anything else
	hub := connectToHub("hub", *mqttPort)
	defer hub.Disconnect()

	go sendHeartbeat(hub, "hub/1hz") // one message per second, on the second

	// open the persistent data store
	glog.Infoln("opening data store:", *dataStore)
	options := bolt.Options{Timeout: time.Second}
	db, err := bolt.Open(*dataStore, 0600, &options)
	if err != nil {
		glog.Fatalln("db:", err)
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

func connectToHub(clientName, hubPort string) *service.Client {
	var err error

	// retry a few times, the internal MQTT server may still be starting up
	for i := 0; i < 3; i++ {
		msg := message.NewConnectMessage()
		msg.SetVersion(4)
		msg.SetCleanSession(true)
		msg.SetClientId([]byte(clientName))
		//msg.SetKeepAlive(10)
		msg.SetWillQos(1)
		msg.SetWillTopic([]byte("will"))
		msg.SetWillMessage([]byte("send me home"))

		client := &service.Client{}
		err = client.Connect("tcp://"+hubPort, msg)
		if err == nil {
			glog.Debugln("connected:", clientName, hubPort)
			return client
		}

		glog.Debugln("cannot connect to MQTT, retrying", err)
		time.Sleep(time.Second)
	}

	glog.Fatal(err)
	return nil
}

func topicAsEvents(hub *service.Client, pattern string) chan Event {
	feed := make(chan Event)

	msg := message.NewSubscribeMessage()
	msg.AddTopic([]byte(pattern), 0)
	e := hub.Subscribe(msg, nil, func(msg *message.PublishMessage) error {
		feed <- Event{
			topic:   string(msg.Topic()),
			payload: msg.Payload(),
		}
		return nil
	})
	if e != nil {
		glog.Fatal(e)
	}

	return feed
}

func startMqttServer(port string) {
	srv := service.Server{}

	certFile := os.Getenv("HUB_MQTT_CERT")
	keyFile := os.Getenv("HUB_MQTT_KEY")

	if certFile != "" && keyFile != "" {
		glog.Infoln("starting MQTT server at", port)
		// TODO: fix surgemq, see https://github.com/surgemq/surgemq/issues/8
		//glog.Fatal(srv.ListenAndServeTLS("tcp://" + port, certFile, keyFile))
		glog.Fatal(srv.ListenAndServe("tcp://" + port))
	} else {
		glog.Infoln("starting MQTT server at", port)
		glog.Fatal(srv.ListenAndServe("tcp://" + port))
	}
}

func startHttpServer(port string) {
	http.HandleFunc("/bar",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", r.URL.Path)
		})

	certFile := os.Getenv("HUB_HTTP_CERT")
	keyFile := os.Getenv("HUB_HTTP_KEY")

	if certFile != "" && keyFile != "" {
		glog.Infoln("starting HTTPS (TLS) server at", port)
		glog.Fatal(http.ListenAndServeTLS(port, certFile, keyFile, nil))
	} else {
		glog.Infoln("starting HTTP server at", port)
		glog.Fatal(http.ListenAndServe(port, nil))
	}
}

func sendHeartbeat(hub *service.Client, topic string) {
	for {
		time.Sleep(time.Duration(1e9 - time.Now().UnixNano()%1e9))

		// publish the heartbeat msg if it's within 10ms of the second mark
		nanos := time.Now().UnixNano()
		if nanos%1e9 < 10e6 {
			millis := strconv.FormatInt(nanos/1e6, 10)

			msg := message.NewPublishMessage()
			msg.SetTopic([]byte(topic))
			msg.SetPayload([]byte(millis))
			e := hub.Publish(msg, func(m, a message.Message, err error) error {
				return nil
			})
			if e != nil {
				glog.Error(e)
				// FIXME re-open new connection to hub after error, yikes!
				hub.Disconnect()
				hub = connectToHub("hub2", *mqttPort)
			}
		} else {
			glog.Errorln("missed heartbeat:", nanos)
		}
		glog.Infoln("heartbeat:", nanos)
	}
}
