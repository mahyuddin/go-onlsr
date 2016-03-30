package main

import (
//  "fmt"
//	"time"

	"github.com/go-ndn/log"
//	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
//	"github.com/go-ndn/tlv"
)

// Neighbour struct to store captured neighbour
type Neighbour struct {
	Address string
    Network string
    Cost uint64
	LocalFace *face
    RemoteFace *face
}


func newNeighbour()() {

}

// function to create link between local and available neighbour
func createLink (availableNeighbour struct{
                    Network, Address string
                    Cost uint64
                }, linkedNeighbours map[string]Neighbour) {
    
    // only add new available neighbour as remote face if not linked
    if _, linked := linkedNeighbours[availableNeighbour.Address]; !linked {

        // create interest channel
        interestChan := make(chan *ndn.Interest)
        
        // remote face
        remote, err := newFace(availableNeighbour.Network, availableNeighbour.Address, availableNeighbour.Cost, interestChan)
        if err != nil {
            delete(linkedNeighbours, availableNeighbour.Address)
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
            linkedNeighbours[availableNeighbour.Address]  = Neighbour {
                
                    Network : availableNeighbour.Network, 
                    Address : availableNeighbour.Address,
                    Cost : availableNeighbour.Cost,
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
    
}