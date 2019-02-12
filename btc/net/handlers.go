package net

import (
	"github.com/chainswatch/bclib/parser"
	"github.com/chainswatch/bclib/serial"

	"fmt"
	log "github.com/sirupsen/logrus"
)

// HandleObject manages and saves tx and block messages
func (p *Peer) HandleObject(object string, payload []byte) (*Inv, error) {
	var hash [32]byte
	copy(hash[:], serial.DoubleSha256(payload))
	inventory := &Inv{
		object:    object,
		raw:       payload,
		timestamp: 0,
		fromIP:    nil,
	}
	if err := p.queue.Push(hash, inventory); err != nil {
		return nil, err
	}
	return inventory, nil
}

//sendheaders
//sendcmpct

// HandlePing replies pong to ping
func (p *Peer) handlePing(nonce []byte) error {
	return p.SendPong(nonce)
}

//feefilter

// HandleAddr parse peer addresses (version >= 31402)
func (p *Peer) handleAddr(payload []byte) ([]*Peer, error) {
	return parseAddr(payload)
}

// HandleVersion handles version message
func (p *Peer) HandleVersion(payload []byte) error {
	buf, err := parser.New(payload)
	if err != nil {
		return err
	}
	p.version = buf.ReadUint32()
	p.services = buf.ReadUint64()
	p.timestamp = buf.ReadUint64()
	return nil
}

// HandleInv parse inventories
func (p *Peer) HandleInv(payload []byte) error {
	inventory, count, err := parseInv(payload)
	if err != nil {
		return err
	}
	return p.sendGetData(inventory, count)
}

// HandleReject prints reject error message
func (p *Peer) HandleReject(payload []byte) {
	log.Warn(fmt.Sprintf("reject (%s): %s", p.ip, payload))
}
