package main

import (
	"log"
	"net"
	"sync"

	"github.com/ugorji/go/codec"
)

func main() {
	log.Println("[server] start")

	l, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	var wg sync.WaitGroup
	wg.Add(2) // stop once client requests have been handled

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}

			go func(c net.Conn) {
				defer c.Close()
				defer wg.Done()

				h := &codec.MsgpackHandle{RawToString: true}
				dec := codec.NewDecoder(c, h)

				for {
					var v interface{}
					if dec.Decode(&v) != nil {
						break
					}

					log.Println("got:", v)
				}
			}(conn)
		}
	}()

	wg.Wait()
	log.Println("[server] end")
}
