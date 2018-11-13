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

type tx struct {
	timestamp			uint32
	fromIP				net.IP
	raw						[]byte
}

type Peer struct {
	ip						net.IP
	port					uint16
	timestamp			uint32
	services			[]byte

	rw						*bufio.ReadWriter
	txs						map[[32]byte]tx
	nextTxs				map[[32]byte]bool		// Buffer waiting for data
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
