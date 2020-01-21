package peering

import (
	"encoding/gob"
	"log"
	"net"
)

type Handler func(*gob.Decoder) bool

func handleConnections(conn net.Conn, handler Handler) {
	dec := gob.NewDecoder(conn)
	defer conn.Close()
	for {
		quit := handler(dec)
		if quit {
			return
		}
	}
}

func Listen(hostport string, handler Handler) error {
	listener, err := net.Listen("tcp", hostport)
	if err != nil {
		return err
	}
	log.Println("Listen on", listener.Addr().String())
	for {
		log.Println("Waiting for connections")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection:", err)
			continue
		}
		go handleConnections(conn, handler)
	}
}
