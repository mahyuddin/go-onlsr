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
type neighbour struct {
	Address    string
	Network    string
	Cost       uint64
	LocalFace  *face
	RemoteFace *face
}

func newNeighbour() {

}

// function to create link between local and available neighbour
func createLink(neighbourNode struct {
	Network, Address string
	Cost             uint64
}, linkedNeighbours map[string]neighbour) {

	// only add unlinked neighbour as remote face
	if _, linked := linkedNeighbours[neighbourNode.Address]; !linked {

		// create interest channel
		interestChan := make(chan *ndn.Interest)

		// remote face
		remote, err := newFace(neighbourNode.Network, neighbourNode.Address, neighbourNode.Cost, interestChan)
		if err != nil {
			delete(linkedNeighbours, neighbourNode.Address)
			log.Println(err)
			//remote.Close()
		} else {
			// local face
			local, err := newFace(config.Local.Network, config.Local.Address, 0, nil)
			if err != nil {
				log.Fatalln(err)
			}
			defer local.Close()

			defer remote.Close()

			// Register remote face as linked neighbour
			linkedNeighbours[neighbourNode.Address] =
				neighbour{
					Network:    neighbourNode.Network,
					Address:    neighbourNode.Address,
					Cost:       neighbourNode.Cost,
					LocalFace:  local,
					RemoteFace: remote,
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
