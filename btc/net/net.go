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

// Inv structure
type Inv struct {
	object    string
	timestamp uint32
	fromIP    net.IP
	raw       []byte
}

// Peer holds information about connected peer
type Peer struct {
	ip        net.IP
	port      uint16
	timestamp uint64
	version		uint32
	services  uint64

	rw       *bufio.ReadWriter
	queue    *Queue // Stores raw txs and blocks
}

// Network holds information about the network status
type Network struct {
	version   uint32
	services  uint32
	userAgent string
	port      uint32

	peers     map[string]*Peer
	maxPeers	uint32
}

// Message holds components of a network message
type Message struct {
	cmd     string
	length  uint32
	payload []byte
}
