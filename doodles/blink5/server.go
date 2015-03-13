package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/surgemq/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func main() {
	srv := &service.Server{}

	enabled := false

	go func() {
		blink := false

		for range time.Tick(500 * time.Millisecond) {
			if enabled {
				blink = !blink

				data, err := json.Marshal(blink)
				if err != nil {
					log.Fatal(err)
				}

				pubmsg := message.NewPublishMessage()
				pubmsg.SetTopic([]byte("blink"))
				pubmsg.SetPayload(data)

				if err := srv.Publish(pubmsg, nil); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	go func() {
		time.Sleep(time.Second) // TODO: hack, give server time to start up

		clt := &service.Client{}

		msg := message.NewConnectMessage()
		msg.SetVersion(4)
		msg.SetClientId([]byte("server"))

		if err := clt.Connect("tcp://:1883", msg); err != nil {
			log.Fatal(err)
		}

		submsg := message.NewSubscribeMessage()
		submsg.AddTopic([]byte("enabled"), 0)

		clt.Subscribe(submsg, nil, func(msg *message.PublishMessage) error {
			err := json.Unmarshal(msg.Payload(), &enabled)
			return err
		})
	}()

	log.Fatal(srv.ListenAndServe("tcp://:1883"))
}
