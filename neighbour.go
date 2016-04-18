package main

import (
	//  "fmt"
	"encoding/hex"
	"net"
	"strings"
	"time"

	"github.com/go-ndn/log"
	//	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	//	"github.com/go-ndn/tlv"
	"github.com/go-ndn/packet"
)

const maxDatagramSize = 8192

var linkedNeighbours map[string]neighbour

// Neighbour struct to store captured neighbour
type neighbour struct {
	Address    string
	Network    string
	Cost       uint64
	LocalFace  *face
	RemoteFace *face
}

type remoteNode struct {
	Network, Address string
	Cost             uint64
}

func multicastMsgHandler(dataSource *net.UDPAddr, numOfBytes int, dataBytes []byte) (node remoteNode) {
	// temporary logic. For neighbourhood discovery, we will replace it to other logic
	if !selfIPAddress(dataSource.IP.String()) {
		if _, linked := linkedNeighbours[dataSource.IP.String()]; !linked {
			address := dataSource.IP.String() + ":6363"
			node = remoteNode{"udp", address, 0}
			log.Println(numOfBytes, "bytes read from", dataSource)
			log.Println(hex.Dump(dataBytes[:numOfBytes]))
		} else {
			log.Println("linked!")
		}

	}
	return
}

func selfIPAddress(recvIPAddress string) (selfIP bool) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		ipAddr := strings.Split(addr.String(), "/")
		if recvIPAddress == ipAddr[0] {
			selfIP = true
		}
	}
	return
}

func serveMulticastUDP(tempNode chan<- remoteNode) {
	addr, err := net.ResolveUDPAddr(config.Multicast.Network, config.Multicast.Address)
	if err != nil {
		log.Fatalln("ResolveUDPAddr error:", err)
	}
	l, err := net.ListenMulticastUDP(config.Multicast.Network, nil, addr)
	if err != nil {
		log.Fatalln("ListenMulticastUDP error:", err)
	}
	err = l.SetReadBuffer(maxDatagramSize)
	if err != nil {
		log.Fatalln("SetReadBuffer error:", err)
	}
	for {
		dataBytes := make([]byte, maxDatagramSize)
		numOfBytes, dataSource, err := l.ReadFromUDP(dataBytes)
		if err != nil {
			log.Fatalln("ReadFromUDP failed:", err)
		}
		node := multicastMsgHandler(dataSource, numOfBytes, dataBytes)
		tempNode <- node
	}
}

func sendUDPHelloPacket() {
	srvAddr := config.Multicast.Address
	addr, err := net.ResolveUDPAddr(config.Multicast.Network, srvAddr)
	if err != nil {
		log.Fatalln(err)
	}
	connect, err := net.DialUDP(config.Multicast.Network, nil, addr)
	for {
		connect.Write([]byte("hello\n"))
		time.Sleep(time.Duration(config.HelloInterval) * time.Second)
	}
}

func neighbourhoodDiscovery(neighbourChan chan<- remoteNode) {
	log.Println("Neighbourhood discovery...")
	tempNode := make(chan remoteNode)
	go serveMulticastUDP(tempNode)
	go sendUDPHelloPacket()
	for {
		node := <-tempNode
		neighbourChan <- node
		time.Sleep(time.Duration(config.HelloInterval) * time.Second)
	}

}

func checkLinkedNeighbour() {

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
	log.Println("Linked neighbours :", len(linkedNeighbours))
}

// function to create link between local and available neighbour
func createLink(neighbourChan <-chan remoteNode) {
	neighbourNode := <-neighbourChan

	// only add unlinked neighbour as remote face
	if _, linked := linkedNeighbours[neighbourNode.Address]; !linked {

		log.Println("Create link between remote and local face")

		// create interest channel
		interestChan := make(chan *ndn.Interest)

		// remote face
		remote, err := newFace(neighbourNode.Network, neighbourNode.Address, neighbourNode.Cost, interestChan)
		if err != nil {
			delete(linkedNeighbours, neighbourNode.Address)
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
