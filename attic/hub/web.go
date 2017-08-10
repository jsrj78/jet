package main

import (
	"log"
)

/*
	// start up the built-in HTTP server
	if *httpPort != "" {
		go func() {
			defer close(quit)
			startHTTPServer(*httpPort)
		}()
	}
*/

func webListener(feed string) {
	for evt := range topicWatcher(feed) {
		log.Println("web:", evt.Topic, "value:", string(evt.Payload))
	}
}
