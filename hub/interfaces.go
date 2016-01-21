package main

import (
    "bufio"
    "fmt"

    "github.com/chimera/rs232"
    "github.com/surge/glog"
)

const SERIAL_DEVICE_PREFIX = "/dev/cu."
//const SERIAL_DEVICE_PREFIX = "/dev/tty"

func listenToSerialPort(device string, baud uint32) {
    options := rs232.Options{ BitRate: baud, DataBits: 8, StopBits: 1 }
    serial, err := rs232.Open(SERIAL_DEVICE_PREFIX + device, options)
    if err != nil {
        glog.Fatal(err)
    }
    scanner := bufio.NewScanner(serial)
    go func() {
        for scanner.Scan() {
            fmt.Println("got:", scanner.Text())
        }
        glog.Fatal("unexpected EOF", serial)
    }()
}
