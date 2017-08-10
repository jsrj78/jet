package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"
)

// hexToVarint interprets the packet data as hex-code varints
func hexToVarint(msg string) []int {
	var values []int
	if hexData, err := hex.DecodeString(msg); err == nil {
		v := 0
		for _, b := range hexData {
			v = (v << 7) | int(b&0x7F)
			if (b & 0x80) != 0 {
				values = append(values, v)
				v = 0
			}
		}
	}
	return values
}

// adminCmd dispatches "jet <cmd> ..." command-line requests.
func adminCmd(decode bool) {
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
			msg := string(evt.Payload)
			fmt.Println(evt.Topic, "=", msg)
			if !decode {
				continue
			}
			if seg := strings.Split(msg, " "); len(seg) == 3 {
				if strings.ToUpper(seg[0]) == "RF69" && len(seg[1]) == 20 {
					if b, err := hex.DecodeString(seg[1]); err == nil {
						// : rf69-info ( -- )  \ display params as hex string
						//   rf69.freq @ h.4 rf69.group @ h.2
						//	 rf.rssi @ h.2 rf.lna @ h.2 rf.afc @ h.4 ;
						freq := int(b[0])*256 + int(b[1])
						afc := int(b[5])*256 + int(b[4])
						if afc >= 32768 {
							afc -= 65536
						}
						rssi := float32(b[3]) * -0.5
						fmt.Printf("    f: %d g: %d rssi: %g lna: %d afc: %d",
							freq, b[2], rssi, b[4], afc)
						fmt.Printf(" dst: %d src: %d hdr: %d len: %d\n",
							b[7]&0x3F, b[8]&0x3F, b[8]>>6, b[9])
					}
				}
				values := hexToVarint(seg[2])
				if len(values) > 0 {
					fmt.Println("       ", values)
				}
			}
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
