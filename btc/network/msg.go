package network

import (
	"git.posc.in/cw/watchers/serial"

	"math/rand"
	"encoding/binary"
	"bytes"
	"time"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type msg struct {
	cmd				string
	length		uint32
	payload		[]byte
}

// SendRawMsg sends command and payload
func (n *Network) sendMsg(pid uint32, cmd string, pl []byte) error {
	var sbuf [24]byte

	binary.LittleEndian.PutUint32(sbuf[0:4], n.networkMagic)
	copy(sbuf[4:16], cmd) // version
	binary.LittleEndian.PutUint32(sbuf[16:20], uint32(len(pl)))

	chksum := serial.DoubleSha256(pl[:])
	copy(sbuf[20:24], chksum[:4])

	msg := append(sbuf[:], pl...)

	p := n.peers[pid]
	log.Info(fmt.Sprintf("Sending %x", msg))
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

func checkType(received string, expected string) error {
	if received != expected {
		return fmt.Errorf("checkType: Unexpected response from peer. Received %s != %s", received, expected)
	}
	return nil
}

//sendheaders
//sendcmpct
//ping
//pong
func (n *Network) msgPong(nonce []byte) {
	n.sendMsg(0, "pong", nonce) // TODO: Replace 0
}

//addr
//feefilter
//inv

// NetworkVersion sends the protocol version to the selected peer and check its response
func (n *Network) msgVersion(id uint32) (*msg, error) {
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

	err := n.sendMsg(id, "version", b.Bytes())
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
func (n *Network) msgVerack(id uint32) (*msg, error) {
	peer := n.peers[id]
	err := n.sendMsg(id, "verack", nil)
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
