package peering

import (
	"log"
	"net"
)

func Connect(addr string) (*net.Conn, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}
