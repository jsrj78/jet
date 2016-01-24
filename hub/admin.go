package main

import (
	"flag"
	"fmt"
)

// adminCmd dispatches "jet <cmd> ..." command-line requests.
func adminCmd() {
	cmd := flag.Arg(0)
	if cmd == "" {
		fmt.Println("JET v" + version)
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

		sendToHub(cmdFlags.Arg(0), []byte(cmdFlags.Arg(1)), *retain)

	case "sub":
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 1 {
			fmt.Println("Usage: jet sub <topic>")
			return
		}

		for evt := range topicWatcher(cmdFlags.Arg(0)) {
			fmt.Printf("%s = %v\n", evt.Topic, evt.Payload)
		}

	case "delete": // unregister a "stuck" registration, i.e. a missing will
		cmdFlags.Parse(cmdArgs)
		if cmdFlags.NArg() != 1 {
			fmt.Println("Usage: jet unreg <topic>")
		}

		sendToHub(cmdFlags.Arg(0), []byte{}, true)

	case "test":
		cmdFlags.Parse(cmdArgs)

		sendToHub("abc", make([]byte, 1024), false)
	}
}
