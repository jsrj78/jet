// This package is a small demo "pack" which connects to JET/Hub.
package main

import (
	"flag"

	"github.com/dataence/glog"
	"github.com/jeelabs/jet/hub/connect"
)

func main() {
	flag.Parse()

	conn, err := connect.NewConnection("demo")
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("connected %q", conn)
	<-conn.Done
	glog.Infof("disconnected %q", conn)
}
