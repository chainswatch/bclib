package network

import (
	"math/rand"
	"encoding/binary"
	"bytes"
	"time"
	"fmt"
	"net"
)

// New initializes peer structure
func (p *Peer) New() {
	p.ip = net.ParseIP("37.59.38.74")
	p.port = "8333"
}

// New initializes network structure
func (n *Network) New() {
	n.networkMagic = 0xD9B4BEF9 // Maybe LE
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333
	n.nPeers = 0
}

func (n *Network) AddPeer(p Peer) {
	n.peers = append(n.peers, p)
	n.nPeers++
}

func (n *Network) netAddr() {
}

// NetworkVersion sends the protocol version to the selected peer
func (n *Network) NetworkVersion(id uint32) ([]byte, error) {
	if id >= n.nPeers {
		return nil, fmt.Errorf("NetworkVersion: (id %d) >= (nPeers %d)", id, n.nPeers)
	}
	peer := n.peers[id]

	b := bytes.NewBuffer([]byte{})

	binary.Write(b, binary.LittleEndian, uint32(n.version)) // Protocol version, 70015
	binary.Write(b, binary.LittleEndian, uint64(n.services)) // Network services
	binary.Write(b, binary.LittleEndian, uint64(time.Now().Unix())) // Timestamp

	// Network address of receiver (26)
	// b.Write(c.PeerAddr.NetAddr.Bytes())
	b.Write(peer.ip) // Network address of receiver
	b.Write([]byte(peer.port)) // Network port of receiver

	// Network address of emitter (26)
	b.Write(bytes.Repeat([]byte{0}, 26))

	binary.Write(b, binary.LittleEndian, uint64(rand.Intn(2^64))) // nonce, 8 bytes

	binary.Write(b, binary.LittleEndian, uint64(len(n.userAgent)))
	b.Write([]byte(n.userAgent))
	// b.Write([]byte{0})

	binary.Write(b, binary.LittleEndian, uint32(0)) // Last blockheight received
	b.WriteByte(1)  // don't notify me about txs (BIP37)

	response, err := n.networkMsg(id, "version", b.Bytes())
	if err != nil {
		return nil, err
	}
	return response, nil
}
