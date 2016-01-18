package main

import (
    "flag"
    "fmt"
    "net/http"
    "time"

    "github.com/boltdb/bolt"
    "github.com/surge/glog"
    "github.com/surgemq/surgemq/service"
)

func main () {
    flag.Parse()

    fmt.Println("JET-Hub")

    // open backing store
    options := bolt.Options{Timeout: time.Second}
    db, err := bolt.Open("storage.db", 0600, &options)
    if err != nil {
        glog.Fatal(err)
    }
    defer db.Close()

    // launch HTTP server
    go func() {
        http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Hello, %q", r.URL.Path)
        })
        glog.Info("Web server ready at http://127.0.0.1:8080")
        glog.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
    }()

    // launch MQTT server
    srv := service.Server{}
    glog.Fatal(srv.ListenAndServe("tcp://127.0.0.1:1883"))
}
