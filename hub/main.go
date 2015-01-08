package main

import "fmt"

func main() {
	// these variables are bumped/updated by goxc when running "make dist"
	var VERSION = "0.0.9-alpha"
	var SOURCE_DATE = "2015-01-08T23:42:28+01:00"

	fmt.Printf("JET/Hub %s (%.10s)\n", VERSION, SOURCE_DATE)
}
