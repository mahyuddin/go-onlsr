package main

import (
	// "fmt" // temporary
	"encoding/json"
	"flag"
	"os"
<<<<<<< HEAD
	"time"
=======
	"sync"
>>>>>>> origin/master

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
<<<<<<< HEAD

	// get neighbour list from config file
	neighbourList := config.Remote

	linkedNeighbours := make(map[string]Neighbour)

    // dummy loop for temporary solution.
	for neighbourList != nil {
        
        log.Println("Neighbour size :", len(linkedNeighbours))

		// Check if linked neighbour is still available.
        // Temporarily, we use packet.Dial to check remote node is running ndn forwareder.
        // We will figure out a better approaches.
		for address, neighbour := range linkedNeighbours {
            checkLink, err := packet.Dial(neighbour.Network, address)
            if err != nil {
                delete(linkedNeighbours, address)
                log.Println(err)
            } else {
                checkLink.Close()
            }
		}
        
		for _, neighbour := range neighbourList {
			go func(network string, address string, cost uint64) {

				// only add new neighbour as remote face if not linked
                if _, linked := linkedNeighbours[address]; !linked {

                    // create interest channel
                    interestChan := make(chan *ndn.Interest)
                    
                    // remote face
                    remote, err := newFace(network, address, cost, interestChan)
                    if err != nil {
                        delete(linkedNeighbours, address)
                        log.Println(err)
                    } else {               
                        // local face
                        local, err := newFace(config.Local.Network, config.Local.Address, 0, nil)
                        if err != nil {
                            log.Fatalln(err)
                        }
                        defer local.Close()
                        
                        defer remote.Close()
                        
                        // Register remote face as linked neighbour
                        linkedNeighbours[address]  = Neighbour {
                            
                                Network : network, 
                                Address : address,
                                Cost : cost,
                                LocalFace : local, 
                                RemoteFace : remote,
                        }    

                        // advertise name prefix
                        go local.advertise(remote)

                        // create remote tunnel
                        for interest := range interestChan {
                            local.ServeNDN(remote, interest)
                        }
                    }
                    
                }
			}(neighbour.Network, neighbour.Address, neighbour.Cost)
		}

		helloIntv := time.Duration(config.HelloInterval) * time.Second
		time.Sleep(helloIntv)
	} // while loop

=======
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

>>>>>>> origin/master
}
