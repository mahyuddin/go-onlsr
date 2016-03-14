package main

import (
//	"time"

	"github.com/go-ndn/log"
//	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
//	"github.com/go-ndn/tlv"
)

type Neighbour struct {
	Address string
	Cost uint64
	ndn.Face
	log.Logger
}

func NewNeighbour()() {
	
}