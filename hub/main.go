package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/surge/glog"
	"github.com/surgemq/surgemq/service"
)

var defFlag = flag.String("def", "hub.def", "hub definition filename")
var runFlag = flag.String("run", ".", "run directory path")

// these variables are bumped/updated by goxc when running "make dist"
var VERSION = "0.0.12-alpha"
var SOURCE_DATE = "2015-01-12T02:01:02+01:00"

func printVersion() {
	fmt.Printf("JET/Hub %s (%.10s)\n", VERSION, SOURCE_DATE)
}

func main() {
	flag.Parse()

	if err := os.MkdirAll(*runFlag, 0755); err != nil {
		glog.Fatal(err)
	}
	if err := os.Chdir(*runFlag); err != nil {
		glog.Fatal(err)
	}

	daemonSetup()
}

func dispatch(cmd string) {
	// the quit, stop, and reload commands have already been handled
	switch cmd {
	case "start":
		daemonStart()
	case "version":
		printVersion()
	case "help":
		flag.PrintDefaults()
	default:
		glog.Fatalln("no such command:", cmd)
	}
}

var srv = &service.Server{}

func worker() {
	if err := srv.ListenAndServe("tcp://:1883"); err != nil {
		glog.Fatal(err)
	}
}

func terminate(wait bool) {
	if wait {
		srv.Close()
	}
}

func reload() {
	glog.Infoln("configuration reloaded")
}
