package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/ugorji/go/codec"
)

var debug = flag.Bool("d", false, "debug mode")

func main() {
	flag.Parse()

	if *debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	log.Print("hubbie1", os.Args)

	launchHubServer()
}

func launchHubServer() {
	addr := config("ADDR", "127.0.0.1")
	port := config("PORT", "7777")
	log.Print("listening on ", addr, ":", port)

	l, err := net.Listen("tcp", addr+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		go func(c net.Conn) {
			defer c.Close()
			defer log.Println("disconnect", c.RemoteAddr())

			log.Println("connect", c.RemoteAddr())
			h := &codec.MsgpackHandle{RawToString: true}
			dec := codec.NewDecoder(c, h)

			for {
				var v interface{}
				if dec.Decode(&v) == nil {
					break
				}
				dispatch(v)
			}
		}(conn)
	}
}

func config(name, defval string) string {
	s := os.Getenv(name)
	if len(s) == 0 {
		s = defval
	}
	return s
}

func dispatch(v interface{}) {
	log.Println("got:", v)
}
