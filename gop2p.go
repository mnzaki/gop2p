package main

import (
	"encoding/gob"
	"log"
	"os"
	"fmt"

	"github.com/mnzaki/gop2p/crdt"
	"github.com/mnzaki/gop2p/peering"
)

var myCounter crdt.GCounter

func main() {
	var myId uint
	fmt.Sscanf(os.Args[1], "%v", &myId)

	log.SetFlags(log.Ltime)
	hostport := os.Args[2]
	peers := os.Args[3:]
	log.Printf("Peers: %v", peers)

	myCounter = crdt.MakeGCounter(crdt.ID(myId))
	myCounter.Increment()
	for _, peer := range peers {
		send, err := peering.Connect(peer)
		if err != nil {
			log.Printf("Error connecting: %v", err)
		} else {
			err := send(myCounter)
			if err != nil {
				log.Printf("Error while sending: %v", err)
			} else {
				log.Printf("Sent %v to %v", myCounter, peer)
			}
		}
	}

	err := peering.Listen(hostport, handleGCounter)
	if err != nil {
		log.Println(err)
	}
}

func handleGCounter(decoder *gob.Decoder) bool {
	var g crdt.GCounter
	err := decoder.Decode(&g)
	if err != nil {
		log.Printf("decode error: %v", err)
		return true
	}

	log.Printf("MyCounter: %v", myCounter)
	myCounter = myCounter.Merge(&g)
	log.Printf("Received: %v", g)
	log.Printf("MyCounter Updated: %v", myCounter)
	return false
}

