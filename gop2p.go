package main

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/mnzaki/gop2p/crdt"
	"github.com/mnzaki/gop2p/peering"
)

func main() {
	log.SetFlags(log.Ltime)
	hostport := os.Args[1]
	peers := os.Args[2:]

	log.Printf("Peers: %v", peers)
	err := peering.Listen(hostport, handleGCounter)
	if err != nil {
		log.Println(err)
	}
}

func handleGCounter(decoder *gob.Decoder) bool {
	var g crdt.GCounter
	err := decoder.Decode(&g)
	if err != nil {
		log.Fatal("decode error:", err)
		return true
	}

	log.Printf("Received: %v", g)
	return false
}

