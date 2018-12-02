package net

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
	invs					map[[32]byte]*inv		// Stores raw txs and blocks
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

// AddPeer adds a new peer
func (n *Network) AddPeer(ip string, port uint16) error {
	peer := Peer{}
	if err := peer.newConnection(ip, port); err != nil {
		return err
	}
	n.peers = append(n.peers, peer)
	n.nPeers++
	return peer.handshake(n.version, n.services, n.userAgent)
}

// New initializes network structure
func (n *Network) New() {
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333
	n.nPeers = 0
}

// apply is passed as an argument to Watch
type apply func(*Peer, *Message) error

// Watch connected peers and apply fn when a message is received
func (n *Network) Watch(fn apply) error {
	// TODO: process peers in parallel (or one by one with select?)
	peer := n.peers[0]
	for {
		msg, err := peer.waitMsg()
		if err != nil {
			return err
		}
		if err = fn(&peer, msg); err != nil {
			return err
		}
	}
	return nil
}
