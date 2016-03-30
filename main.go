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
	neighbourList := config.Remote

	linkedNeighbours := make(map[string]Neighbour)

    // dummy loop for temporary solution.
	for neighbourList != nil {
        
        log.Println("Neighbour size :", len(linkedNeighbours))

		// Check if linked neighbour is still available.
		for address, neighbour := range linkedNeighbours {
            checkLink, err := packet.Dial(neighbour.Network, address)
            if err != nil {
                delete(linkedNeighbours, address)
                log.Println(err)
            } else {
                checkLink.Close()
            }
		}
        
        // create links between local forwarder and available neighbour forwarders
		for _, neighbour := range neighbourList {
            go createLink(neighbour, linkedNeighbours)
		}

        // dummy hello interval
		helloIntv := time.Duration(config.HelloInterval) * time.Second
		time.Sleep(helloIntv)
	} // while loop
}