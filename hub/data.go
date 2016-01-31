package main

import (
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

// dataStoreInit initialises the data store, creating it if necessary
func dataStoreInit(file string) *bolt.DB {
	log.Println("opening data store:", file)
	options := bolt.Options{Timeout: time.Second}
	var e error
	db, e = bolt.Open(file, 0600, &options)
	if e != nil {
		log.Fatalln("data:", e)
	}
	return db
}

// dataModifyListener listens for data store and delete requests
func dataModifyListener(feed string) {
	for evt := range topicWatcher(feed) {
		key := strings.Split(evt.Topic, "/")[1:]
		if len(key) == 0 || key[0] == "" {
			log.Println("modify key missing, value:", len(evt.Payload), "bytes")
			continue
		}

		if len(evt.Payload) > 0 {
			storeValue(key, evt.Payload)
		} else {
			deleteKey(key)
		}
	}
}

func storeValue(key []string, value []byte) {
	log.Println("store key:", key, "value:", len(value), "bytes")
}

func deleteKey(key []string) {
	log.Println("delete key:", key)
}

// dataAccessListener listens for data fetch and list requests
func dataAccessListener(feed string) {
	for evt := range topicWatcher(feed) {
		key := strings.Split(evt.Topic, "/")[1:]
		if len(key) == 0 || key[0] == "" {
			log.Println("access key missing:", string(evt.Payload))
			continue
		}

		var req struct {
			Cmd, Reply string
		}
		if evt.Decode(&req) {
			if req.Reply == "" {
				log.Println("no reply topic:", req.Cmd, "key:", key)
				continue
			}
			switch req.Cmd {
			case "", "fetch":
				fetchKey(key, req.Reply)
			case "*", "list":
				listKeys(key, req.Reply)
			default:
				log.Println("bad data request:", req.Cmd, "key:", key)
			}
		}
	}
}

func fetchKey(key []string, reply string) {
	log.Println("fetch key:", key, "to:", reply)
}

func listKeys(key []string, reply string) {
	log.Println("list keys:", key, "to:", reply)
}
