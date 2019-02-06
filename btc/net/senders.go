package net

import (
	"github.com/chainswatch/bclib/parser"
	"github.com/chainswatch/bclib/serial"

	"bytes"
	"encoding/binary"
	"math/rand"

	"fmt"
	"time"
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

	// log.Trace(fmt.Sprintf("Sending %s [%x] %x", cmd, sbuf, pl))
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

// SendPing sends a ping message
func (p *Peer) SendPing() error {
	nonce := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonce, uint64(rand.Intn(2^64)))
	return p.sendMsg("ping", nonce)
}

// SendPong sends a pong message to conneted peer
func (p *Peer) SendPong(nonce []byte) error {
	return p.sendMsg("pong", nonce) // TODO: Replace 0
}

// SendHeaders send a sendheaders message to connected peer
func (p *Peer) SendHeaders() error {
	return p.sendMsg("sendheaders", []byte{0})
}

// SendGetdata requests a single block or transaction by hash
// to connected peer
// TODO: Separate them (tx, block) by type when checking uniqueness?
func (p *Peer) sendGetData(inventory [][]byte, count uint64) error {
	b := bytes.NewBuffer([]byte{})
	var hash [32]byte
	newCount := uint64(0)
	for i := uint64(0); i < count; i++ {
		copy(hash[:], inventory[i][4:])
		if !p.queue.Exists(hash) {
			b.Write(inventory[i])
			newCount++
			// log.Debug(fmt.Sprintf("Debug inv: %x", inventory[i][:4]))
		}
	}
	buf := b.Bytes()
	msg := append(parser.CompactSize(newCount), buf...)
	return p.sendMsg("getdata", msg)
}

// sendGetblocks
// sendGetheaders
// sendGetaddr

// sendVersion sends the protocol version to the selected peer and check its response
func (p *Peer) sendVersion(version, services uint32, userAgent string) (*Message, error) {
	b := bytes.NewBuffer([]byte{})

	binary.Write(b, binary.LittleEndian, uint32(version))           // Protocol version, 70015
	binary.Write(b, binary.LittleEndian, uint64(services))          // Network services
	binary.Write(b, binary.LittleEndian, uint64(time.Now().Unix())) // Timestamp

	// Network address of receiver (26)
	b.Write(p.ip)                              // Network address of receiver
	b.Write([]byte(fmt.Sprintf("%d", p.port))) // Network port of receiver

	// Network address of emitter (26)
	b.Write(bytes.Repeat([]byte{0}, 26))

	binary.Write(b, binary.LittleEndian, uint64(rand.Intn(2^64))) // nonce, 8 bytes
	binary.Write(b, binary.LittleEndian, uint64(len(userAgent)))
	b.Write([]byte(userAgent))

	// Last blockheight received
	binary.Write(b, binary.LittleEndian, uint32(0))
	b.WriteByte(1) // don't notify me about txs (BIP37)

	err := p.sendMsg("version", b.Bytes())
	if err != nil {
		return nil, err
	}
	response, err := p.waitMsg()
	if err != nil {
		return nil, err
	}
	err = checkType(response.cmd, "version")
	if err != nil {
		return nil, err
	}
	return response, nil
}

//
func (p *Peer) sendVerack() (*Message, error) {
	err := p.sendMsg("verack", nil)
	if err != nil {
		return nil, err
	}
	response, err := p.waitMsg()
	if err != nil {
		return nil, err
	}
	err = checkType(response.cmd, "verack")
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Handshake sends the Version message (wait for response) followed by a verack message (wait for response)
func (p *Peer) handshake(version, services uint32, userAgent string) error {
	m, err := p.sendVersion(version, services, userAgent)
	if err != nil {
		return err
	}
	// log.Trace(fmt.Sprintf("Received: %s %d %x", m.Cmd(), m.Length(), m.Payload()))
	p.HandleVersion(m.Payload())
	// TODO: Print host name, IP and version

	m, err = p.sendVerack()
	if err != nil {
		return err
	}
	// log.Trace(fmt.Sprintf("Received: %s %d %x", m.Cmd(), m.Length(), m.Payload()))

	p.SendHeaders()
	return nil
}
