package main

import (
	"log"
	"net"

	"github.com/ugorji/go/codec"
)

func main() {
	log.Println("[client] start")

	conn, err := net.Dial("tcp", ":4444")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	h := &codec.MsgpackHandle{}
	enc := codec.NewEncoder(conn, h)

	enc.Encode("abc")
	enc.Encode(123)
	enc.Encode([]string{"A", "B", "C"})

	log.Println("[client] end")
}
