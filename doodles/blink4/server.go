package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/ws", websocket.Handler(wsHandler))

	log.Println("listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func wsHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}
