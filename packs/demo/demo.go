package main

import (
	"flag"

	"github.com/dataence/glog"
	"github.com/jeelabs/jet/hub/connect"
)

func main() {
	flag.Parse()

	conn, err := connect.NewConnection("ha")
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("conn %v", conn)

	done := make(chan struct{})
	<-done
}
