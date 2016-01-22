package main

import (
	"flag"
	"fmt"

	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

var adminQuit = make(chan struct{})

func adminCmd(hub *service.Client) {
	cmd := flag.Arg(0)
	cmdArgs := flag.Args()[1:]
	cmdFlags := flag.NewFlagSet(cmd, flag.ExitOnError)

	switch cmd {

	default:
		fmt.Println("ha!")
		close(adminQuit)

	case "pub":
		retain := cmdFlags.Bool("r", false, "send with RETAIN flag set")
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 2 {
			fmt.Println("Usage: jet pub ?-r? topic payload")
			return
		}

		pubmsg := message.NewPublishMessage()
		pubmsg.SetTopic([]byte(cmdFlags.Arg(0)))
		pubmsg.SetPayload([]byte(cmdFlags.Arg(1)))
		pubmsg.SetRetain(*retain)
		pubmsg.SetQoS(1)
		hub.Publish(pubmsg, adminDone)

	case "sub":
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 1 {
			fmt.Println("Usage: jet sub topic")
			return
		}

		submsg := message.NewSubscribeMessage()
		submsg.AddTopic([]byte(cmdFlags.Arg(0)), 0)
		hub.Subscribe(submsg, nil, func(msg *message.PublishMessage) error {
			fmt.Printf("%s = %q\n", msg.Topic(), msg.Payload())
			return nil
		})

	case "test":
		cmdFlags.Parse(cmdArgs)

		submsg := message.NewSubscribeMessage()
		submsg.AddTopic([]byte("abc"), 0)
		hub.Subscribe(submsg, nil, nil)

		pubmsg := message.NewPublishMessage()
		pubmsg.SetTopic([]byte("abc"))
		pubmsg.SetPayload(make([]byte, 1024))
		pubmsg.SetQoS(0)
		hub.Publish(pubmsg, adminDone)
	}

	<-adminQuit // wait forever
}

func adminDone(msg, ack message.Message, err error) error {
	close(adminQuit)
	return nil
}
