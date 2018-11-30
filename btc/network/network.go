package network

import (
	"net"
	"bufio"
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
	object				string
	timestamp			uint32
	fromIP				net.IP
	raw						[]byte
}

// Peer holds information about connected peer
type Peer struct {
	ip						net.IP
	port					uint16
	timestamp			uint32
	services			[]byte

	rw						*bufio.ReadWriter
	invs					map[[32]byte]inv		// Stores raw txs and blocks
	nextInvs			map[[32]byte]bool		// Buffer waiting for data
}

// Network holds information about the network status
type Network struct {
	version				uint32
	services			uint32
	userAgent			string
	port					uint32
	peers					[]Peer
	nPeers				uint32
}

type msg struct {
	cmd				string
	length		uint32
	payload		[]byte
}
