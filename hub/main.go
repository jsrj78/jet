package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/boltdb/bolt"
    "github.com/surgemq/surgemq/service"
)

func main () {
    flag.Parse()

    log.Print("[JET/Hub]")

    // open backing store
    log.Print("opening database: storage.db")
    options := bolt.Options{Timeout: time.Second}
    db, err := bolt.Open("storage.db", 0600, &options)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // launch HTTP server
    go func() {
        http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Hello, %q", r.URL.Path)
        })
        log.Print("Web server ready at http://127.0.0.1:8080")
        log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
    }()

    // launch MQTT server
    srv := service.Server{}
    log.Print("MQTT server ready at tcp://127.0.0.1:1883")
    log.Fatal(srv.ListenAndServe("tcp://127.0.0.1:1883"))
}
