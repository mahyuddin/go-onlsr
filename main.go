package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/go-ndn/log"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/packet"
)

var (
	flagConfig = flag.String("config", "go-onlsr.json", "config path")
	flagDebug  = flag.Bool("debug", false, "enable logging")
)

var (
	key ndn.Key
)

func main() {
	flag.Parse()

	// config
	configFile, err := os.Open(*flagConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	// key
	pem, err := os.Open(config.PrivateKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer pem.Close()

	key, err = ndn.DecodePrivateKey(pem)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("key", key.Locator())

	// experiment using dummy neighbour list

	// get neighbour list from config file
	neighboursList := config.Remote

	// Create linkedNeighbour struct map
	// to store list of linked neighbours
	linkedNeighbours := make(map[string]neighbour)

	// dummy loop for temporary solution.
	for neighboursList != nil {

		log.Println("Neighbourhood size :", len(linkedNeighbours))

		// Check if linked neighbour is still available.
		for address, linkedNode := range linkedNeighbours {
			// cuba cari cara lain utk connection checking
			if checkLink, err := packet.Dial(linkedNode.Network, address); err != nil {
				if linkedNode.RemoteFace.Handler == nil {
					log.Println("Face tarak ada...")
				}
				delete(linkedNeighbours, address)
				log.Println(err)
			} else {
				checkLink.Close()
			}
		}

		// create links between local forwarder and available neighbour forwarders
		for _, availableNode := range neighboursList {
			go createLink(availableNode, linkedNeighbours)
		}

		// dummy hello interval
		helloIntv := time.Duration(config.HelloInterval) * time.Second
		time.Sleep(helloIntv)
	} // dummy while loop
}
