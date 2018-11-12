package network

import (
	"math/rand"
	"encoding/binary"
	"bytes"
	"time"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func messageType(msg []byte) string {
	return fmt.Sprintf("%s", bytes.Trim(msg[:12], "\x00"))
}

func checkType(msg []byte, expected string) error {
	received := messageType(msg)
	log.Debug("Received:", received)
	if received != expected {
		return fmt.Errorf("checkType: Unexpected response from peer. Received %s != %s", received, expected)
	}
	return nil
}

//sendheaders
//sendcmpct
//ping
//addr
//feefilter
//inv

// NetworkVersion sends the protocol version to the selected peer
func (n *Network) msgVersion(id uint32) ([]byte, error) {
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

	response, err := n.sendMsg(id, "version", b.Bytes())
	if err != nil {
		return nil, err
	}
	if err = checkType(response, "version"); err != nil {
		log.Warn(err)
		return nil, err
	}
	return response, nil
}

//
func (n *Network) msgVerack(id uint32) ([]byte, error) {
	response, err := n.sendMsg(id, "verack", nil)
	if err != nil {
		return nil, err
	}
	if err = checkType(response, "verack"); err != nil {
		log.Warn(err)
		return nil, err
	}
	// TODO: Check if response is verack
	return response, nil
}
