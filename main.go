package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/go-ndn/log"
	"github.com/go-ndn/ndn"
	// "github.com/go-ndn/packet"
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

	// --------------- experiment using dummy neighbour list ---------------------

	neighbourChan := make(chan remoteNode)
	go neighbourhoodDiscovery(neighbourChan)

	// dummy loop for temporary solution.
	for {

		go checkLinkedNeighbour()
		go createLink(neighbourChan)

		// dummy hello interval
		time.Sleep(time.Duration(config.HelloInterval) * time.Second)
	} // dummy for loop
}
