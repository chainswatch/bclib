package net

import (
	"git.posc.in/cw/bclib/parser"
	"git.posc.in/cw/bclib/serial"

	"math/rand"
	"encoding/binary"
	"bytes"

	"time"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func checkType(received string, expected string) error {
	if received != expected {
		return fmt.Errorf("checkType: Unexpected response from peer. Received %s != %s", received, expected)
	}
	return nil
}

// sendMsg sends command and payload
func (p *Peer) sendMsg(cmd string, pl []byte) error {
	var sbuf [24]byte

	binary.LittleEndian.PutUint32(sbuf[0:4], networkMagic)
	copy(sbuf[4:16], cmd) // version
	binary.LittleEndian.PutUint32(sbuf[16:20], uint32(len(pl)))

	chksum := serial.DoubleSha256(pl[:])
	copy(sbuf[20:24], chksum[:4])

	msg := append(sbuf[:], pl...)

	log.Info(fmt.Sprintf("Sending [%x] %x", sbuf, pl))
	_, err := p.rw.Write(msg)
	if err != nil {
		return err
	}
	err = p.rw.Flush()
	if err != nil {
		return err
	}
	return nil
}

// sendPong sends a pong message to conneted peer
func (p *Peer) SendPong(nonce []byte) {
	p.sendMsg("pong", nonce) // TODO: Replace 0
}

// sendGetdata requests a single block or transaction by hash
// to connected peer
func (p *Peer) SendGetdata(inventory [][]byte, count uint64) {
	b := bytes.NewBuffer([]byte{})
	b.Write(parser.Varint(count))
	var hash [32]byte
	for i := uint64(0); i < count; i++ {
		copy(hash[:], inventory[i][4:])
		if _, found := p.invs[hash]; !found { // if not exist
			p.nextInvs[hash] = true
			b.Write(inventory[i])
		}
	}
	p.sendMsg("getdata", b.Bytes())
}

// sendGetblocks
// sendGetheaders
// sendGetaddr

// sendVersion sends the protocol version to the selected peer and check its response
func (n *Network) sendVersion(id uint32) (*Message, error) {
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
	b.Write([]byte(fmt.Sprintf("%d", peer.port))) // Network port of receiver

	// Network address of emitter (26)
	b.Write(bytes.Repeat([]byte{0}, 26))

	binary.Write(b, binary.LittleEndian, uint64(rand.Intn(2^64))) // nonce, 8 bytes
	binary.Write(b, binary.LittleEndian, uint64(len(n.userAgent)))
	b.Write([]byte(n.userAgent))

	// Last blockheight received
	binary.Write(b, binary.LittleEndian, uint32(0))
	b.WriteByte(1)	// don't notify me about txs (BIP37)

	err := peer.sendMsg("version", b.Bytes())
	if err != nil {
		return nil, err
	}
	response, err := peer.WaitMsg()
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
func (n *Network) sendVerack(id uint32) (*Message, error) {
	peer := n.peers[id]
	err := peer.sendMsg("verack", nil)
	if err != nil {
		return nil, err
	}
	response, err := peer.WaitMsg()
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

// Handshake sends the Version message (wait for response) followed by a verack message (wait for response)
func (n *Network) handshake(peerID uint32) error {
	response, err := n.sendVersion(peerID)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("response %x", response))

	response, err = n.sendVerack(peerID)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("response %x", response))
	return nil
}
