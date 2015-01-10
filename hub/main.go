package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var hubDefFlag = flag.String("d", "hub.def", "hub definition filename")
var runDirFlag = flag.String("r", ".", "run directory path")

// these variables are bumped/updated by goxc when running "make dist"
var VERSION = "0.0.9-alpha"
var SOURCE_DATE = "2015-01-08T23:42:28+01:00"

func printVersion() {
	fmt.Printf("JET/Hub %s (%.10s)\n", VERSION, SOURCE_DATE)
}

func main() {
	flag.Parse()
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
		log.Fatalln("no such command:", cmd)
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
	log.Println("terminating...")
	stop <- struct{}{}
	if wait {
		<-done
	}
}

func reload() {
	log.Println("configuration reloaded")
}
