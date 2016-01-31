package main

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
		log.Println("serial:", serName)

		var firstByte uint8
		if len(evt.Payload) > 0 {
			firstByte = evt.Payload[0]
		}

		switch firstByte {
			case 0, '{':
				if port, ok := portMap[serName]; ok {
					log.Println("serial:", serName, "closing")
					port.Close()
					delete(portMap, serName)
				}
		}

		if firstByte == '{' {
			var serReq struct {
				Device, SendTo string
			}
			if evt.Decode(&serReq) {
				serial := listenToSerial(serName, serReq.Device, serReq.SendTo)
				if serial != nil {
					portMap[serName] = serial
				} else {
					delete(portMap, serName)
				}
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
