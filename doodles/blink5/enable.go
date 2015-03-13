package main

import (
	"log"
	"os"
	"time"

	"github.com/surgemq/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func main() {
	clt := &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetClientId([]byte("enabler"))

	if err := clt.Connect("tcp://:1883", msg); err != nil {
		log.Fatal(err)
	}

	pubmsg := message.NewPublishMessage()
	pubmsg.SetTopic([]byte("enabled"))
	pubmsg.SetPayload([]byte(os.Args[1]))

	err := clt.Publish(pubmsg, nil)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)
}
