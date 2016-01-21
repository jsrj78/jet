package main

import (
    "log"
    "os"

    "github.com/surgemq/message"
    "github.com/surgemq/surgemq/service"
)

func adminCmd(hub *service.Client) {
    log.SetFlags(log.Lshortfile)
    log.Println("starting:", os.Args[1:])

    submsg := message.NewSubscribeMessage()
    submsg.AddTopic([]byte("abc"), 0)

    hub.Subscribe(submsg, nil, nil)

    pubmsg := message.NewPublishMessage()
    pubmsg.SetTopic([]byte("abc"))
    pubmsg.SetPayload(make([]byte, 1024))
    pubmsg.SetQoS(0)

    hub.Publish(pubmsg, nil)
}
