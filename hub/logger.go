package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func timestampRepeater(feed chan event) {
	for evt := range feed {
		millis := time.Now().UnixNano() / 1e6
		topic := fmt.Sprintf("%s/%d", evt.Topic, millis)
		publish(topic, evt.Payload, false)
	}
}

func loggerSaveToDisk(dir string, feed chan event) {
	var lastPath string
	var lastFile *os.File

	for evt := range feed {
		// topic = "logger/<device>/<milliseconds>"
		segments := strings.Split(evt.Topic, "/")

		var millis int64
		if _, e := fmt.Sscanf(segments[2], "%d", &millis); e != nil {
			log.Println("scan error:", evt.Topic, e)
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
		// linefeeds must be escaped, these log files are one-entry-per-line
		msg := strings.Replace(string(evt.Payload), "\n", "\\n", -1)
		fmt.Fprintf(lastFile, "L %s %s %s\n", tod, segments[1], msg)
	}
}
