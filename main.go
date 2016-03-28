package main

import (
	// "fmt" // temporary
	"encoding/json"
	"flag"
	"os"
	"sync"

	"github.com/go-ndn/log"
	"github.com/go-ndn/ndn"
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
	var waitGroup sync.WaitGroup

    // get neighbour list from config file
	neighbourList := config.Remote

	for _, neighbour := range neighbourList {

        go func(network string, address string, cost uint64) {
            
            // local face
            local, err := newFace(config.Local.Network, config.Local.Address, 0, nil)
            if err != nil {
                log.Fatalln(err)
            }
            defer local.Close()

            // create interest channel
            interestChan := make(chan *ndn.Interest)

            // remote face
            remote, err := newFace(network, address, cost, interestChan)
            if err != nil {
                log.Fatalln(err)
            }
            defer remote.Close()

            // advertise name prefix
            go local.advertise(remote)

            // create remote tunnel
            for interest := range interestChan {
                local.ServeNDN(remote, interest, waitGroup)
            }

            waitGroup.Done()

        }(neighbour.Network, neighbour.Address, neighbour.Cost)
	}
    
	waitGroup.Add(len(neighbourList))
	waitGroup.Wait()

}
