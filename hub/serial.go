package main

import (
	"bufio"
	"log"

	"github.com/chimera/rs232"
)

const SERIAL_DEVICE_PREFIX = "/dev/cu."

//const SERIAL_DEVICE_PREFIX = "/dev/tty"

func listenToDevices(changes chan Event) {
	listenToSerialPort("usbserial-A40119DV", 57600)
	//listenToSerialPort("USB0", 57600)
}

func listenToSerialPort(device string, baud uint32) {
	options := rs232.Options{BitRate: baud, DataBits: 8, StopBits: 1}
	serial, err := rs232.Open(SERIAL_DEVICE_PREFIX+device, options)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(serial)
	go func() {
		for scanner.Scan() {
			log.Println("got:", scanner.Text())
		}
		log.Fatal("unexpected EOF", serial)
	}()
}
