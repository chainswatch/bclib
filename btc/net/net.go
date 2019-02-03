package net

import (
	"bufio"
	"net"
)

/*
func NetworkMessage() {
	// network magic
	// command
	// payload length
	// payload checksum
	// payload
}
*/

const (
	networkMagic = 0xD9B4BEF9
)

// Inventory structure
type inv struct {
	object    string
	timestamp uint32
	fromIP    net.IP
	raw       []byte
}

// Peer holds information about connected peer
type Peer struct {
	ip        net.IP
	port      uint16
	timestamp uint32
	services  []byte

	rw       *bufio.ReadWriter
	invs     map[[32]byte]*inv // Stores raw txs and blocks
	nextInvs map[[32]byte]bool // Buffer waiting for data
}

// Network holds information about the network status
type Network struct {
	version   uint32
	services  uint32
	userAgent string
	port      uint32

	peers     []Peer
	nPeers    uint32
	maxPeers	uint32
}

// Message holds components of a network message
type Message struct {
	cmd     string
	length  uint32
	payload []byte
}

//
func (n *Network) ConnectedPeers() uint32 {
	return n.nPeers
}

// New initializes network structure
func (n *Network) New() {
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333
	n.nPeers = 0
	n.maxPeers = 10
}

// apply is passed as an argument to Watch
type apply func(*Peer, *Message, interface{}) error

// Watch connected peers and apply fn when a message is received
func (n *Network) Watch(fn apply, argFn interface{}) {
	for _,p := range n.peers {
		p := p
		go p.handleConnection(fn, argFn)
	}
}

