package main

import "fmt"

func main() {
	// these variables are bumped/updated by goxc when running "make dist"
	var VERSION = "0.0.7-alpha"
	var SOURCE_DATE = "2015-01-08T19:08:37+01:00"

	fmt.Printf("JET/Server %s (%.10s)\n", VERSION, SOURCE_DATE)
}
