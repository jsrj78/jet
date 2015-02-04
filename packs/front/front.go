// Front end web server and web socket support.
package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/jeelabs/jet/hub/connect"
	"github.com/surge/glog"
	"golang.org/x/net/websocket"
)

func main() {
	flag.Parse()

	conn, err := connect.NewConnection("front")
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("connected %v", conn)

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
	glog.Infoln("ws", ws)
	io.Copy(ws, ws) // just echo for now
}
