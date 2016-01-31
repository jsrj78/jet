package main

// TODO refactor to use a struct with member methods for each serial connection

import (
	"bufio"
	"log"

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
				}
				if evt.Decode(&req) {
					ser := listenToSerial(serName, req.Device, req.SendTo)
					if ser != nil {
						portMap[serName] = ser
					}
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
func listenToSerial(name, device, topic string) *rs232.Port {
	options := rs232.Options{BitRate: 57600, DataBits: 8, StopBits: 1}
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
			// TODO support binary data with Bytes() i.s.o. Text() as option?
			//	this is likely to break in some places, e.g. "jet sub ..."
			// solution could be to use some special topic naming convention
			//	also would need to enable read timeout for framing the data
			// sendToHub(topic, scanner.Bytes(), false)
			sendToHub(topic, scanner.Text(), false)
		}
		log.Println("serial:", name, "EOF")
	}()
	return serial
}

// processSerialRequests interprets a list of special serial port requests.
func processSerialRequests(name string, port *rs232.Port, reqs []interface{}) {
	log.Println("serial", name, "requests:", reqs)

	for _, req := range reqs {
		log.Printf("req: %T %v\n", req, req)
		/*
			case '1'..'9':
				var req uint32
				if evt.Decode(&req) {
					log.Println("serial", serName, "baudrate:", req)
				}
		*/
	}
}
