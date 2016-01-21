package main

import (
    "bufio"
    "flag"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/boltdb/bolt"
    "github.com/chimera/rs232"
    "github.com/surge/glog"
    "github.com/surgemq/surgemq/service"
)

func main () {
    if len(os.Args) > 1 && os.Args[1] == "admin" {
        admin()
        return
    }

    flag.Parse()
    glog.Info(append([]string{"JET/Hub"}, os.Args[1:]...))

    // open backing store
    glog.Info("opening database: storage.db")
    options := bolt.Options{Timeout: time.Second}
    db, err := bolt.Open("storage.db", 0600, &options)
    if err != nil {
        glog.Fatal(err)
    }
    defer db.Close()

    // open serial port
    listenToSerialPort("/dev/tty.usbserial-A40119DV", 57600)

    // launch HTTP server
    go func() {
        http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Hello, %q", r.URL.Path)
        })
        glog.Infoln("starting HTTP server at", ":8947")
        glog.Fatal(http.ListenAndServe(":8947", nil))
    }()

    // launch MQTT server
    srv := service.Server{}
    glog.Infoln("starting MQTT server at", ":1883")
    glog.Fatal(srv.ListenAndServe("tcp://" + ":1883"))
}

func listenToSerialPort(device string, baud uint32) {
    options := rs232.Options{ BitRate: baud, DataBits: 8, StopBits: 1 }
    serial, err := rs232.Open(device, options)
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
