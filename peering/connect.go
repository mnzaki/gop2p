package peering

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
)

type Sender func(v interface{}) error

func Connect(addr string, handler Handler) (Sender, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	enc := gob.NewEncoder(rw)
	go func() {
		dec := gob.NewDecoder(rw)
		for {
			quit := handler(dec)
			if quit {
				return
			}
		}
	}()

	return func(v interface{}) error {
		defer rw.Flush()
		return enc.Encode(v)
	}, nil
}
