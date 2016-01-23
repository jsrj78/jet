package main

import (
	"flag"
	"log"
	"fmt"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

func adminCmd(hub *mqtt.Client) {
	cmd := flag.Arg(0)
	if cmd == "" {
		fmt.Println("JET v" + VERSION)
		return
	}
	cmdArgs := flag.Args()[1:]
	cmdFlags := flag.NewFlagSet(cmd, flag.ExitOnError)

	switch cmd {

	default:
		fmt.Println("Available commands: pub sub delete test")

	case "pub":
		retain := cmdFlags.Bool("r", false, "send with RETAIN flag set")
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 2 {
			fmt.Println("Usage: jet pub ?-r? <topic> <payload>")
			return
		}

		t := hub.Publish(cmdFlags.Arg(0), 0, *retain, cmdFlags.Arg(1))
		if t.Wait() && t.Error() != nil {
			log.Fatal(t.Error())
		}

	case "sub":
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 1 {
			fmt.Println("Usage: jet sub <topic>")
			return
		}

		topic := cmdFlags.Arg(0)
		t := hub.Subscribe(topic, 0, func(hub *mqtt.Client, msg mqtt.Message) {
			fmt.Printf("%s = %q\n", msg.Topic(), msg.Payload())
		})
		if t.Wait() && t.Error() != nil {
			log.Fatal(t.Error())
		}

		quit := make(chan struct{})
		<-quit // this waits forever

	case "delete": // unregister a "stuck" registration, i.e. a missing will
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 1 {
			fmt.Println("Usage: jet unreg <topic>")
		}

		t := hub.Publish(cmdFlags.Arg(0), 1, true, []byte{})
		if t.Wait() && t.Error() != nil {
			log.Fatal(t.Error())
		}

	case "test":
		cmdFlags.Parse(cmdArgs)

		t := hub.Publish("abc", 0, false, make([]byte, 1024))
		if t.Wait() && t.Error() != nil {
			log.Fatal(t.Error())
		}
	}
}
