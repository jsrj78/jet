package main

import (
	"flag"

	"github.com/dataence/glog"
	"github.com/jeelabs/jet/hub/connect"
)

var done = make(chan struct{})

func main() {
	flag.Parse()

	conn, err := connect.NewConnection("demo")
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("connected %q", conn)
	<-done
	glog.Infof("disconnected %q", conn)
}
