package net

import (
	"bufio"
	"net"

	zmq "github.com/pebbe/zmq4"
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
	payload   []byte
	hash      [32]byte
}

// Peer holds information about connected peer
type Peer struct {
	ip        net.IP
	port      uint16
	timestamp uint64
	version   uint32
	services  uint64
	errors    uint16

	rw    *bufio.ReadWriter
	conn  net.Conn
	Pub   *zmq.Socket
	queue *Queue // Stores raw txs and blocks
}

// apply is passed as an argument to Watch
type apply func(*Peer, *Message) error

// Network holds information about the network status
type Network struct {
	version   uint32
	services  uint32
	userAgent string
	port      uint32

	peers    map[string]*Peer // Connected peers
	newAddr  map[string]*Peer // New peers to connect to
	banned   map[string]bool
	maxPeers uint32

	fn apply
}

// Message holds components of a network message
type Message struct {
	cmd     string
	length  uint32
	payload []byte
}
