package main

// TODO refactor to use a struct with member methods for each serial connection

import (
	"bufio"
	"log"
	"time"

	"github.com/chimera/rs232"
)

// serialProcessRequests handles all serial port setup and outgoing data.
func serialProcessRequests(feed string) {
	portMap := map[string]*rs232.Port{}

	for evt := range topicWatcher(feed) {
		serName := evt.Topic[7:] // TODO wrong if feed isn't "serial/+"

		if len(evt.Payload) == 0 || evt.Payload[0] == '{' {
			if port, ok := portMap[serName]; ok {
				log.Println("serial:", serName, "closing")
				port.Close()
				delete(portMap, serName)
			}
		}

		if len(evt.Payload) > 0 {
			port, _ := portMap[serName]
			switch evt.Payload[0] {
			case '{':
				var req struct {
					Device, SendTo string
					Baud           uint32
					Init           []interface{}
				}
				if evt.Decode(&req) {
					ser := listenToSerial(serName, req.Device, req.SendTo,
						req.Baud)
					if ser != nil {
						portMap[serName] = ser
					}
					processSerialRequests(serName, ser, req.Init)
				}
			case '[':
				var req []interface{}
				if evt.Decode(&req) && port != nil {
					processSerialRequests(serName, port, req)
				}
			case '"':
				var req string
				if evt.Decode(&req) && port != nil {
					port.Write([]byte(req))
				}
			default:
				n := len(evt.Payload)
				log.Println("serial", serName, "ignored:", n, "bytes")
			}
		}
	}
}

// listenToSerial reads incoming serial text lines and publishes them to MQTT.
func listenToSerial(name, device, topic string, baud uint32) *rs232.Port {
	if baud == 0 {
		baud = 57600
	}
	options := rs232.Options{BitRate: baud, DataBits: 8, StopBits: 1}
	serial, err := rs232.Open(device, options)
	if err != nil {
		log.Print(err)
		return nil
	}

	// TODO flush old pending data, since it's not actually real-time
	//if n, _ := serial.BytesAvailable(); n > 0 {
	//	serial.Read(make([]byte, n))
	//}

	scanner := bufio.NewScanner(serial)
	go func() {
		for scanner.Scan() {
			sendToHub(topic, scanner.Bytes(), false)
		}
		log.Println("serial:", name, "EOF")
	}()
	return serial
}

// processSerialRequests interprets a list of special serial port requests.
func processSerialRequests(name string, port *rs232.Port, reqs []interface{}) {
	log.Println("serial", name, "requests:", reqs)

	for _, req := range reqs {
		if cmd, ok := req.(string); ok {
			switch cmd {
			case "+dtr":
				port.SetDTR(true)
			case "-dtr":
				port.SetDTR(false)
			case "+rts":
				port.SetRTS(true)
			case "-rts":
				port.SetRTS(false)
			default:
				if len(cmd) > 0 && cmd[0] == '=' {
					port.Write([]byte(cmd[1:]))
				} else {
					log.Println("serial", name, "cmd:", cmd, "?")
				}
			}
		} else if delay, ok := req.(float64); ok {
			if 1 <= delay && delay <= 10000 {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			} else {
				log.Println("serial", name, "delay:", delay, "?")
			}
		} else {
			log.Println("serial", name, "req:", req, "?")
		}
	}
}
