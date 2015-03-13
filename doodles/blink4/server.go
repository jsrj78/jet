package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/ws", websocket.Handler(wsHandler))

	log.Println("listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func wsHandler(ws *websocket.Conn) {
	log.Println("connect")
	defer log.Println("disconnect")

	dec := json.NewDecoder(ws)
	enc := json.NewEncoder(ws)

	// TODO: needs to be mutex-protected
	var in struct {
		Enabled bool `json:"enabled"`
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		defer log.Println("bye")

		var out struct {
			Blink bool `json:"blink"`
		}

		ticker := time.Tick(500 * time.Millisecond)

		for {
			select {
			case <-ticker:
				if in.Enabled {
					out.Blink = !out.Blink

					err := enc.Encode(&out)
					if err != nil {
						log.Fatal(err)
					}
				}
			case <-done:
				return
			}
		}
	}()

	for {
		err := dec.Decode(&in)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println("got", in)
	}
}
