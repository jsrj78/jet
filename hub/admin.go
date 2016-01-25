package main

import (
	"flag"
	"fmt"
)

// adminCmd dispatches "jet <cmd> ..." command-line requests.
func adminCmd() {
	cmd := flag.Arg(0)
	if cmd == "" {
		fmt.Println("JET " + version)
		return
	}
	cmdArgs := flag.Args()[1:]
	cmdFlags := flag.NewFlagSet(cmd, flag.ExitOnError)

	switch cmd {

	default:
		fmt.Println("Available commands: pub sub delete config test")

	case "pub":
		retain := cmdFlags.Bool("r", false, "send with RETAIN flag set")
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 2 {
			fmt.Println("Usage: jet pub ?-r? <topic> <payload>")
			return
		}

		sendToHub(cmdFlags.Arg(0), []byte(cmdFlags.Arg(1)), *retain)

	case "sub":
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 1 {
			fmt.Println("Usage: jet sub <topic>")
			return
		}

		for evt := range topicWatcher(cmdFlags.Arg(0)) {
			fmt.Printf("%s = %s\n", evt.Topic, evt.Payload)
		}

	case "delete": // unregister a "stuck" registration, i.e. a missing will
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 1 {
			fmt.Println("Usage: jet delete <topic>")
			return
		}

		sendToHub(cmdFlags.Arg(0), []byte{}, true)

	case "config":
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 0 {
			fmt.Println("Usage: jet config")
			return
		}

		// show all the retained state in MQTT, which is always sent first
		// TODO minor bug: this hangs if there is no MQTT activity at all
		for evt := range topicWatcher("#") {
			if !evt.Retained {
				break
			}
			fmt.Printf("%s = %s\n", evt.Topic, evt.Payload)
		}

	case "test":
		cmdFlags.Parse(cmdArgs)

		sendToHub("abc", make([]byte, 1024), false)
	}
}
