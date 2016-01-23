package main

import (
	"bufio"
	"encoding/json"
	"log"

	"github.com/chimera/rs232"
)

func processSerialRequests(feed chan Event) {
	portmap := map[string]*rs232.Port{}

	for req := range feed {
		log.Println("req:", req.Topic)

		var serReq struct {
			Device string `json:"device"`
			SendTo string `json:"sendto"`
		}
		if e := json.Unmarshal(req.Payload, &serReq); e != nil {
			log.Println("serial request parse error:", req, e)
		} else {
			serial := listenToSerialPort(serReq.Device, serReq.SendTo)
			if serial != nil {
				portmap[req.Topic] = serial
			} else {
				delete(portmap, req.Topic)
			}
		}
	}
}

func listenToSerialPort(device, topic string) *rs232.Port {
	options := rs232.Options{BitRate: 57600, DataBits: 8, StopBits: 1}
	serial, err := rs232.Open(device, options)
	if err != nil {
		log.Println("cannot open:", device, err)
		return nil
	}

	scanner := bufio.NewScanner(serial)
	go func() {
		for scanner.Scan() {
			log.Println("got:", scanner.Text())
			publish(topic, scanner.Bytes(), false)
		}
		log.Fatal("unexpected EOF", serial)
	}()
	return serial
}
