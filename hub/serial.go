package main

import (
	"bufio"
	"encoding/json"
	"log"

	"github.com/chimera/rs232"
)

func processSerialRequests(feed chan Event) {
	portmap := map[string]*rs232.Port{}

	for evt := range feed {
		log.Println("evt:", evt.Topic)

		var serReq struct {
			Device string `json:"device"`
			SendTo string `json:"sendto"`
		}
		if e := json.Unmarshal(evt.Payload, &serReq); e != nil {
			log.Println("serial request parse error:", evt, e)
		} else {
			serial := listenToSerialPort(serReq.Device, serReq.SendTo)
			if serial != nil {
				portmap[evt.Topic] = serial
			} else {
				delete(portmap, evt.Topic)
			}
		}
	}
}

func listenToSerialPort(device, topic string) *rs232.Port {
	options := rs232.Options{BitRate: 57600, DataBits: 8, StopBits: 1}
	serial, err := rs232.Open(device, options)
	if err != nil {
		log.Print(err)
		return nil
	}

	scanner := bufio.NewScanner(serial)
	go func() {
		for scanner.Scan() {
			publish(topic, scanner.Bytes(), false)
		}
		log.Println("unexpected EOF:", device)
	}()
	return serial
}
