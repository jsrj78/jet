package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
)

var defFlag = flag.String("def", "hub.def", "hub definition filename")
var runFlag = flag.String("run", ".", "run directory path")

// these variables are bumped/updated by goxc when running "make dist"
var VERSION = "0.0.11-alpha"
var SOURCE_DATE = "2015-01-10T14:10:48+01:00"

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

var stop = make(chan struct{})
var done = make(chan struct{})

func worker() {
	for {
		time.Sleep(time.Second)
		if _, ok := <-stop; ok {
			break
		}
	}
	time.Sleep(3 * time.Second)
	done <- struct{}{}
}

func termHandler(wait bool) {
	stop <- struct{}{}
	if wait {
		<-done
	}
}

func reload() {
	glog.Infoln("configuration reloaded")
}
