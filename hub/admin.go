package main

import (
    "log"
    "os"

    "github.com/surgemq/message"
    "github.com/surgemq/surgemq/service"
)

func admin() {
    log.SetFlags(log.Lshortfile)
    log.Println(os.Args[1:])

    hub := &service.Client{}

    msg := message.NewConnectMessage()
    msg.SetVersion(4)
    msg.SetCleanSession(true)
    msg.SetClientId([]byte("ad" + "min"))
    msg.SetKeepAlive(10)
    msg.SetWillQos(1)
    msg.SetWillTopic([]byte("will"))
    msg.SetWillMessage([]byte("send me home"))

    if err := hub.Connect("tcp://127.0.0.1:1883", msg); err != nil {
        log.Fatal(err)
    }

    submsg := message.NewSubscribeMessage()
    submsg.AddTopic([]byte("abc"), 0)

    hub.Subscribe(submsg, nil, nil)

    pubmsg := message.NewPublishMessage()
    pubmsg.SetTopic([]byte("abc"))
    pubmsg.SetPayload(make([]byte, 1024))
    pubmsg.SetQoS(0)

    log.Println("999")
    hub.Publish(pubmsg, nil)

    log.Println("123")
    //hub.Disconnect()
    log.Println("456")
}
