// Front end web server and web socket support.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/jeelabs/jet/hub/connect"
	"github.com/surge/glog"
	"golang.org/x/net/websocket"
)

var (
	conn *connect.Connection

	muSess sync.RWMutex
	sess   = map[string]*websocket.Conn{}
)

func main() {
	flag.Parse()

	var err error
	conn, err = connect.NewConnection("front")
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("connected %v", conn)

	// subscribe to all topics
	conn.Listen("#", func(key string, val interface{}) {
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
		for _, ws := range sess {
			ws.Write(msg)
		}
	})

	// set up static web server and websocket handler
	http.Handle("/", http.FileServer(http.Dir("../../web")))
	http.Handle("/ws", websocket.Handler(wsHandler))

	glog.Info("listening to :1111")
	log.Fatal(http.ListenAndServe(":1111", nil))

	// never reached
	//<-conn.Done
	//glog.Infof("disconnected %v", conn)
}

func wsHandler(ws *websocket.Conn) {
	remote := ws.Request().RemoteAddr

	conn.Send("hub/connect", remote)
	defer conn.Send("hub/disconnect", remote)

	muSess.Lock()
	sess[remote] = ws
	muSess.Unlock()

	buf := bufio.NewReader(ws)
	pendingPrefix := ""
	var err error

	for {
		// grab until next space as topic name
		var suffix string
		suffix, err = buf.ReadString(' ')
		if err != nil {
			if err == io.EOF {
				return // websocket was closed
			}
			break
		}
		topic := pendingPrefix + suffix[:len(suffix)-1]

		// then decode one JSON object
		decoder := json.NewDecoder(buf)
		var payload interface{}
		err = decoder.Decode(&payload)
		if err != nil {
			break
		}

		// everything we read too far ends up as prefix for the next topic
		var all []byte
		all, err = ioutil.ReadAll(decoder.Buffered())
		if err != nil {
			break
		}
		pendingPrefix = string(all)

		// FIXME what do do if the readahead consumed more than a topic prefix?
		if strings.IndexByte(pendingPrefix, ' ') >= 0 {
			glog.Fatal("JSON decoding was too greedy:", pendingPrefix)
		}

		// publish the topic/payload we got
		conn.Send(topic, payload)
	}

	// can only exit the above loop if something went wrong
	glog.Errorln(remote, err)
}
