package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// loggerTimestamper resends each incoming message to a new timestamped topic.
func loggerTimestamper(feed string) {
	for evt := range topicWatcher(feed) {
		millis := time.Now().UnixNano() / 1e6
		topic := fmt.Sprintf("%s/%d", evt.Topic, millis)
		sendToHub(topic, evt.Payload, false)
	}
}

// loggerSaveToDisk picks up timestamped messages and saves them to log files.
func loggerSaveToDisk(feed, dir string) {
	var lastPath string
	var lastFile *os.File

	for evt := range topicWatcher(feed) {
		var message string
		if !evt.Decode(&message) {
			continue // ignore this event, failure has been logged by Decode
		}
		// linefeeds must be escaped, since log files have one-entry-per-line
		message = strings.Replace(message, "\n", "\\n", -1)

		// topic = "logger/<device>/<milliseconds>"
		segments := strings.Split(evt.Topic, "/")

		// extract the timestamp from the topic
		var millis int64
		if _, e := fmt.Sscanf(segments[2], "%d", &millis); e != nil {
			log.Println("timestamp error:", evt.Topic, e)
			continue
		}
		stamp := time.Unix(0, millis*1e6).UTC()

		// rotate to a new file every day at 0:00 (UTC)
		path := stamp.Format("/2006/20060102.txt")
		if path != lastPath {
			var e error
			if e = os.MkdirAll(dir+path[:5], 0777); e != nil {
				log.Fatal(e)
			}
			if lastFile != nil {
				lastFile.Close()
			}
			fileOpts := os.O_CREATE | os.O_APPEND | os.O_WRONLY
			lastFile, e = os.OpenFile(dir+path, fileOpts, 0666)
			if e != nil {
				log.Fatal(e)
			}
			lastPath = path
			log.Println("logger output to:", dir+path)
		}

		tod := stamp.Format("15:04:05.000")
		fmt.Fprintf(lastFile, "L %s %s %s\n", tod, segments[1], message)
	}
}
