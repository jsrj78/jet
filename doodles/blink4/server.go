package main

import (
	"encoding/json"
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
	dec := json.NewDecoder(ws)
	enc := json.NewEncoder(ws)

	var in struct {
		Enabled bool `json:"enabled"`
	}

	var out struct {
		Blink bool `json:"blink"`
	}

	for {
		err := dec.Decode(&in)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println("got", in)

		out.Blink = in.Enabled

		err = enc.Encode(&out)
		if err != nil {
			log.Fatal(err)
		}
	}
}
