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
var peerMap map[string]peering.Sender

func main() {
	var myId uint
	peerMap = make(map[string]peering.Sender)
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
			peerMap[peer] = send
			err := send(myCounter)
			if err != nil {
				log.Printf("Error while sending: %v", err)
			} else {
				log.Printf("Sent %v to %v", myCounter, peer)
			}
		}
	}

	go func() {
		var cmd string
		for {
			fmt.Scanf("%v", &cmd)
			if cmd == "i" {
				myCounter.Increment()
				log.Printf("Incremented to: %v", myCounter)
				replicate()
			} else if cmd == "p" {
				log.Printf("%v", myCounter)
			}
		}
	}()

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

	oldVal := myCounter.Value()

	log.Printf("MyCounter: %v", myCounter)
	myCounter = myCounter.Merge(&g)
	log.Printf("Received: %v", g)
	log.Printf("MyCounter Updated: %v", myCounter)

	newVal := myCounter.Value()

	if oldVal != newVal {
		replicate()
	}
	return false
}

func replicate() {
	for addr, send := range peerMap {
		log.Printf("replicating to peer %v", addr)
		err := send(myCounter)
		if err != nil {
			log.Printf("replication error: %v", err)
		}
	}
}
