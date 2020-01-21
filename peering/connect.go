package peering

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
)

type Sender func(v interface{}) error

func Connect(addr string) (Sender, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(conn)
	enc := gob.NewEncoder(writer)
	return func(v interface{}) error {
		defer writer.Flush()
		return enc.Encode(v)
	}, nil
}
