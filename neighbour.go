package main

import (
//  "fmt"
//	"time"

//	"github.com/go-ndn/log"
//	"github.com/go-ndn/mux"
//	"github.com/go-ndn/ndn"
//	"github.com/go-ndn/tlv"
)

type neighbour struct {
	Address string
<<<<<<< HEAD
    Network string
    Cost uint64
	LocalFace *face
    RemoteFace *face
}

type AvailableNeighbour struct {
	Address string
    Network string
	MacAddress string
    Cost uint64
=======
	Cost uint64
	Face face
}

type availableNeighbour struct {
	Address string
	MacAddress string
>>>>>>> origin/master
}

func newNeighbour()() {

}