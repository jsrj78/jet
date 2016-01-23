package main

import (
	"fmt"
	"time"
)

func timestampRepeater(feed chan Event) {
	for evt := range feed {
		millis := time.Now().UnixNano() / 1e6
		topic := fmt.Sprintf("%s/%d", evt.Topic, millis)
		publish(topic, evt.Payload, false)
	}
}

func loggerSaveToDisk(dir string, feed chan Event) {
	for evt := range feed {
		_ = evt // TODO
	}
}
