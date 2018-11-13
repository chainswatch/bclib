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

type Peer struct {
	ip						net.IP
	port					uint16
	timestamp			uint32
	services			[]byte

	rw						*bufio.ReadWriter
}

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
