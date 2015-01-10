package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

// these variables are bumped/updated by goxc when running "make dist"
var VERSION = "0.0.9-alpha"
var SOURCE_DATE = "2015-01-08T23:42:28+01:00"

func main() {
	flag.Parse()
	daemonAndSignalSetup()
}

func usage(verbose bool) {
	if verbose {
		fmt.Println(LONG_USAGE)
	} else {
		fmt.Printf("JET/Hub %s (%.10s)\n", VERSION, SOURCE_DATE)
	}
}

var LONG_USAGE = `
  Usage: jethub <cmd> {options...}

    start     - start the server in the background
    quit      - quit the server gently (SIGTERM)
    stop      - stop the server focefully (SIGQUIT)
    reload    - make the server reload its configuration file
    version   - display version information
`

func performCmd(cmd string) {
	// the quit, stop, and reload commands have already been handled
	switch cmd {
	case "start":
		startDaemon()
	case "help":
		usage(true)
	case "version":
		usage(false)
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
	done <- struct{}{}
}

func termHandler(quit bool) {
	log.Println("terminating...")
	stop <- struct{}{}
	if quit {
		<-done
	}
}

func reloadHandler() {
	log.Println("configuration reloaded")
}
