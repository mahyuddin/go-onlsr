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
	Cost uint64
	Face face
}

type availableNeighbour struct {
	Address string
	MacAddress string
}

func newNeighbour()() {

}