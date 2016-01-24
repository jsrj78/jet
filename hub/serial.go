package main

import (
	"bufio"
	"log"

	"github.com/chimera/rs232"
)

// serialProcessRequests handles all serial port setup and outgoing data.
func serialProcessRequests(feed string) {
	portmap := map[string]*rs232.Port{}

	for evt := range topicWatcher(feed) {
		log.Println("evt:", evt.Topic)

		var serReq struct {
			Device, SendTo string
		}
		if evt.Decode(&serReq) {
			serial := listenToSerial(serReq.Device, serReq.SendTo)
			if serial != nil {
				portmap[evt.Topic] = serial
			} else {
				delete(portmap, evt.Topic)
			}
		}
	}
}

// listenToSerial reads incoming serial text lines and publishes them to MQTT.
func listenToSerial(device, topic string) *rs232.Port {
	options := rs232.Options{BitRate: 57600, DataBits: 8, StopBits: 1}
	serial, err := rs232.Open(device, options)
	if err != nil {
		log.Print(err)
		return nil
	}

	scanner := bufio.NewScanner(serial)
	go func() {
		for scanner.Scan() {
			sendToHub(topic, scanner.Text(), false)
		}
		log.Println("unexpected EOF:", device)
		// TODO: serial.Close() ?
	}()
	return serial
}
