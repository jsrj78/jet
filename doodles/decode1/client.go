package main

import (
	"github.com/ugorji/go/codec"
	"log"
	"net"
)

func main() {
	log.Println("[client]")

	conn, err := net.Dial("tcp", ":7777")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	h := &codec.MsgpackHandle{}
	enc := codec.NewEncoder(conn, h)

	enc.Encode("abc")
	enc.Encode(123)
	enc.Encode([]string{"A", "B", "C"})
}
