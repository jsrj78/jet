package main

import (
	"log"
	"strings"
	"time"
	"bytes"

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
		keys := bytes.Split([]byte(evt.Topic), []byte("/"))
		if len(keys) < 3 {
			log.Println("bad modify key:", evt.Topic)
			continue
		}

		if len(evt.Payload) > 0 {
			log.Println("store:", evt.Topic, "value:", len(evt.Payload), "b")
			storeValue(keys[1:], evt.Payload)
		} else {
			log.Println("delete:", evt.Topic)
			deleteKey(keys[1:])
		}
	}
}

func storeValue(keys [][]byte, value []byte) {
	updater := func(tx *bolt.Tx) error {
		bucket, e := tx.CreateBucketIfNotExists([]byte(keys[0]))
		for i, k := range keys {
			if e != nil {
				break
			}
			if len(k) == 0 {
				continue
			}
			if i < len(keys)-1 {
				bucket, e = bucket.CreateBucketIfNotExists(k)
			} else {
				e = bucket.Put(k, value)
			}
		}
		return e
	}
	if e := db.Update(updater); e != nil {
		log.Println("store:", e)
	}
}

func deleteKey(keys [][]byte) {
}

// dataAccessListener listens for data fetch and list requests
func dataAccessListener(feed string) {
	for evt := range topicWatcher(feed) {
		key := strings.Split(evt.Topic, "/")[1:]
		if len(key) < 2 || key[0] == "" {
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
