package network

import (
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

type Peer struct {
	ip						net.IP
	port					string
}

type Network struct {
	networkMagic	uint32
	version				uint32
	services			uint32
	userAgent			string
	port					uint32
	peers					[]Peer
	nPeers				uint32
}
