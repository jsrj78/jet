package main

import (
	"log"

	"github.com/surgemq/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func main() {
	clt := &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetClientId([]byte("dumper"))
	msg.SetKeepAlive(3600)

	if err := clt.Connect("tcp://:1883", msg); err != nil {
		log.Fatal(err)
	}

	submsg := message.NewSubscribeMessage()
	submsg.AddTopic([]byte("#"), 0)

	clt.Subscribe(submsg, nil, func(msg *message.PublishMessage) error {
		log.Println(string(msg.Topic()), string(msg.Payload()))
		return nil
	})

	forever := make(chan struct{})
	<-forever
}
