package main

import (
	"bytes"
	"errors"
	"log"
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
		keys := bytes.Split([]byte(evt.Topic), []byte("/"))
		if len(keys) < 2 {
			log.Println("bad modify key:", evt.Topic)
			continue
		}

		if len(evt.Payload) > 0 {
			log.Println("store:", evt.Topic, "value:", len(evt.Payload), "b")
			storeValue(keys, evt.Payload)
		} else {
			log.Println("delete:", evt.Topic)
			deleteKey(keys)
		}
	}
}

func storeValue(keys [][]byte, value []byte) {
	updater := func(tx *bolt.Tx) error {
		last := len(keys) - 1
		bucket, e := tx.CreateBucketIfNotExists([]byte(keys[1]))
		for i := 2; i < last; i++ {
			if e == nil {
				bucket, e = bucket.CreateBucketIfNotExists(keys[i])
			}
		}
		k := keys[last]
		if e == nil {
			if len(k) > 0 {
				e = bucket.Put(k, value)
			} else {
				// TODO store obj as N items
			}
		}
		return e
	}
	if e := db.Update(updater); e != nil {
		log.Println("store:", e)
	}
}

func deleteKey(keys [][]byte) {
	updater := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(keys[1]))
		last := len(keys) - 1
		for i := 2; i < last; i++ {
			if bucket != nil {
				bucket = bucket.Bucket(keys[i])
			}
		}
		e := errors.New("?")
		if bucket != nil {
			k := keys[last]
			if len(k) > 0 {
				e = bucket.Delete(k)
			} else {
				e = bucket.DeleteBucket(k)
			}
		}
		return e
	}
	if e := db.Update(updater); e != nil {
		log.Println("delete:", e)
	}
}

// dataAccessListener listens for data fetch and list requests
func dataAccessListener(feed string) {
	for evt := range topicWatcher(feed) {
		keys := bytes.Split([]byte(evt.Topic), []byte("/"))
		last := len(keys) - 1
		if last < 1 {
			log.Println("bad access key:", evt.Topic)
			continue
		}

		var reply string
		if evt.Decode(&reply) {
			if reply == "" {
				log.Println("no reply for key:", evt.Topic)
				continue
			}
			if len(keys[last]) > 0 {
				log.Println("fetch key:", evt.Topic, "to:", reply)
				if e := fetchKey(keys, reply); e != nil {
					log.Println("fetch error:", e, "key:", evt.Topic)
				}
			} else {
				log.Println("list keys:", evt.Topic, "to:", reply)
				if e := listKeys(keys, reply); e != nil {
					log.Println("fetch error:", e, "key:", evt.Topic)
				}
			}
		}
	}
}

func fetchKey(keys [][]byte, reply string) error {
	viewer := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(keys[1]))
		last := len(keys) - 1
		for i := 2; i < last; i++ {
			k := keys[i]
			if bucket != nil {
				bucket = bucket.Bucket(k)
			}
		}
		if bucket == nil {
			return errors.New("?")
		}
		sendToHub(reply, bucket.Get(keys[last]), false)
		return nil
	}
	return db.View(viewer)
}

func listKeys(keys [][]byte, reply string) error {
	viewer := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(keys[1]))
		for i := 2; i < len(keys); i++ {
			k := keys[i]
			if bucket != nil {
				bucket = bucket.Bucket(k)
			}
		}
		if bucket == nil {
			return errors.New("?")
		}
		result := map[string]int{}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			result[string(k)] = len(v)
		}
		sendToHub(reply, result, false)
		return nil
	}
	return db.View(viewer)
}
