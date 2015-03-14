// This package is a small demo "pack" which connects to JET/Hub.
package main

import (
	"flag"
	"time"

	"github.com/jeelabs/jet/hub/connect"
	"github.com/surge/glog"
)

func main() {
	flag.Parse()

	conn, err := connect.NewConnection("demo")
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("connected %v", conn)

	// subscribe to all topics
	conn.Listen("#", func(key string, val interface{}) {
		glog.Infof("t: %q p: %v", key, val)
	})

	// send a test message after one second
	go func() {
		time.Sleep(time.Second)
		conn.Send("/test/haha", []interface{}{123, nil, "abc"})
		time.Sleep(time.Second)
		close(conn.Done)
	}()

	<-conn.Done
	glog.Infof("disconnected %v", conn)
}
