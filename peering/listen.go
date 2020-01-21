package peering

import (
	"encoding/gob"
	"log"
	"net"
	"bufio"
)

type Handler func(*gob.Decoder) bool
type NewPeerHandler func(string, Sender)

func handleConnections(conn net.Conn, handler Handler, newPeerHandler NewPeerHandler) {
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	enc := gob.NewEncoder(rw)
	sender := func(v interface{}) error {
		defer rw.Flush()
		return enc.Encode(v)
	}
	newPeerHandler(conn.LocalAddr().String(), sender)

	dec := gob.NewDecoder(rw)
	for {
		quit := handler(dec)
		if quit {
			return
		}
	}
}

func Listen(hostport string, handler Handler, newConnHandler NewPeerHandler) error {
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
		go handleConnections(conn, handler, newConnHandler)
	}
}
