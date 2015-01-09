package main

import (
	"github.com/ugorji/go/codec"
	"log"
	"net"
)

func main() {
	log.Println("[server]")

	l, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			h := &codec.MsgpackHandle{RawToString: true}
			dec := codec.NewDecoder(conn, h)
			for {
				var v interface{}
				if dec.Decode(&v) != nil {
					break
				}
				log.Println(v)
			}
			log.Println("close")
			c.Close()
		}(conn)
	}
}
