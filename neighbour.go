package main

import (
	//  "fmt"
	// "time"

	"github.com/go-ndn/log"
	//	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	//	"github.com/go-ndn/tlv"
	"github.com/go-ndn/packet"
)

// Neighbour struct to store captured neighbour
type neighbour struct {
	Address    string
	Network    string
	Cost       uint64
	LocalFace  *face
	RemoteFace *face
}

func neighbourhoodDiscovery(neighbourChan chan<- []struct {
	Network, Address string
	Cost             uint64
}) {
    // dummy neighbourhood discovery
	log.Println("Neighbourhood discovery...")
	neighbourChan <- config.Remote

}

func checkLinkedNeighbour(linkedNeighbours map[string]neighbour) {

	log.Println("Linked neighbours :", len(linkedNeighbours))

	// Check if linked neighbour is still available.
	for address, linkedNode := range linkedNeighbours {
		// cuba cari cara lain utk connection checking
		if checkLink, err := packet.Dial(linkedNode.Network, address); err != nil {
			if linkedNode.RemoteFace.Handler == nil {
				log.Println("Face not exist")
			}
			delete(linkedNeighbours, address)
			log.Println(err)
		} else {
			checkLink.Close()
		}
	}
}

// function to create link between local and available neighbour
func createLink(neighbourChan <-chan []struct {
	Network, Address string
	Cost             uint64
}, linkedNeighbours map[string]neighbour) {

	log.Println("Create link between remote and local face")

	// linkedNeighbours := <-linkedNodeChan
	neighboursList := <-neighbourChan

	for _, neighbourNode := range neighboursList {

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
}
