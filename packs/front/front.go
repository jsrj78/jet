// Front end web server and web socket support.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/jeelabs/jet/hub/connect"
	"github.com/surge/glog"
	"golang.org/x/net/websocket"
)

var (
	mqttConn *connect.Connection

	muSess     sync.RWMutex
	wsSessions = map[string]*websocket.Conn{}

	persistentTopics = map[string]interface{}{} // all topics matching "/..."
)

func main() {
	flag.Parse()

	var err error
	mqttConn, err = connect.NewConnection("front")
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infoln("connected", mqttConn)

	// subscribe to all topics
	mqttConn.Listen("#", func(key string, val interface{}) {
		updatePersistentTopics(key, val)

		// create a string + JSON message to send out
		data, err := json.Marshal(val)
		if err != nil {
			glog.Errorln(err)
			return
		}
		msg := append([]byte(key+" "), data...)

		// send the message to all web sockets
		muSess.RLock()
		defer muSess.RUnlock()
		for _, ws := range wsSessions {
			ws.Write(msg)
		}
	})

	// set up static web server and websocket handler
	http.Handle("/", http.FileServer(http.Dir("../../web")))
	http.Handle("/ws", websocket.Handler(wsHandler))

	glog.Info("listening to :1111")
	log.Fatal(http.ListenAndServe(":1111", nil))

	// never reached
	//<-mqttConn.Done
	//glog.Infoln("disconnected", mqttConn)
}

func updatePersistentTopics(key string, val interface{}) {
	// keep track of all persistent keys, i.e. topics "/..."
	if strings.HasPrefix(key, "/") {
		// storing nil is treated as a deletion
		if val != nil {
			persistentTopics[key] = val
		} else {
			delete(persistentTopics, key)

			// deleting ".../" deletes all child entries as well
			if strings.HasSuffix(key, "/") {
				for k, _ := range persistentTopics {
					if strings.HasPrefix(k, key) {
						delete(persistentTopics, k)
					}
				}
			}
		}
	}
}

func wsHandler(ws *websocket.Conn) {
	remote := ws.Request().RemoteAddr

	mqttConn.Send("hub/connect", remote)
	defer mqttConn.Send("hub/disconnect", remote)

	muSess.Lock()
	wsSessions[remote] = ws
	muSess.Unlock()

	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				return // websocket was closed
			}
			glog.Errorln(remote, err)
			break
		}

		sep := bytes.IndexByte(buf[:n], ' ')
		if sep < 0 {
			glog.Fatal("no separator:", buf[:n])
		}
		topic := string(buf[:sep])

		var payload interface{}
		err = json.Unmarshal(buf[sep+1:n], &payload)
		if err != nil {
			glog.Error(remote, err)
			break
		}

		// publish the topic/payload we got
		mqttConn.Send(topic, payload)
	}
}
