package network

import (
	"git.posc.in/cw/watchers/parser"

	"math/rand"
	"encoding/binary"
	"bytes"

	"time"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// sendPong sends a pong message to conneted peer
func (p *Peer) sendPong(nonce []byte) {
	p.sendMsg("pong", nonce) // TODO: Replace 0
}

// sendGetdata requests a single block or transaction by hash
// to connected peer
func (p *Peer) sendGetdata(inventory [][]byte, count uint64) {
	b := bytes.NewBuffer([]byte{})

	b.Write(parser.Varint(count))
	for i := uint64(0); i < count; i++ {
		b.Write(inventory[i])
	}
	p.sendMsg("getdata", b.Bytes())
}

// sendGetblocks
// sendGetheaders
// sendGetaddr

// NetworkVersion sends the protocol version to the selected peer and check its response
func (n *Network) sendVersion(id uint32) (*msg, error) {
	if id >= n.nPeers {
		return nil, fmt.Errorf("NetworkVersion: (id %d) >= (nPeers %d)", id, n.nPeers)
	}
	peer := n.peers[id]

	b := bytes.NewBuffer([]byte{})

	binary.Write(b, binary.LittleEndian, uint32(n.version)) // Protocol version, 70015
	binary.Write(b, binary.LittleEndian, uint64(n.services)) // Network services
	binary.Write(b, binary.LittleEndian, uint64(time.Now().Unix())) // Timestamp

	// Network address of receiver (26)
	b.Write(peer.ip) // Network address of receiver
	b.Write([]byte(peer.port)) // Network port of receiver

	// Network address of emitter (26)
	b.Write(bytes.Repeat([]byte{0}, 26))

	binary.Write(b, binary.LittleEndian, uint64(rand.Intn(2^64))) // nonce, 8 bytes
	binary.Write(b, binary.LittleEndian, uint64(len(n.userAgent)))
	b.Write([]byte(n.userAgent))

	binary.Write(b, binary.LittleEndian, uint32(0)) // Last blockheight received
	b.WriteByte(1)  // don't notify me about txs (BIP37)

	err := peer.sendMsg("version", b.Bytes())
	if err != nil {
		return nil, err
	}
	response, err := peer.waitMsg()
	if err != nil {
		return nil, err
	}
	err = checkType(response.cmd, "version")
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	return response, nil
}

//
func (n *Network) sendVerack(id uint32) (*msg, error) {
	peer := n.peers[id]
	err := peer.sendMsg("verack", nil)
	if err != nil {
		return nil, err
	}
	response, err := peer.waitMsg()
	if err != nil {
		return nil, err
	}
	err = checkType(response.cmd, "verack")
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	return response, nil
}
