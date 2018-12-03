package net

import (
  "git.posc.in/cw/bclib/parser"
	"git.posc.in/cw/bclib/serial"

	"fmt"

	log "github.com/sirupsen/logrus"
)

// HandleObject manages tx and block messages
func (p *Peer) HandleObject(object string, payload []byte) error {
	var hash [32]byte
	copy(hash[:], serial.DoubleSha256(payload))
	if _, asked := p.nextInvs[hash]; !asked {
		return fmt.Errorf("handleTx: Hash not found: %x", hash)
	}
	if _, exist := p.invs[hash]; exist {
		return fmt.Errorf("handleTx: %s already exists", object)
	}
	invObj := &inv {
		object: object,
		raw: payload,
		timestamp: 0,
		fromIP: nil,
	}
	p.invs[hash] = invObj
	return nil
}

//sendheaders
//sendcmpct

// HandlePing replies pong to ping
func (p *Peer) HandlePing(nonce []byte) error {
	return p.SendPong(nonce)
}

//feefilter

// HandleAddr parse peer addresses (version >= 31402)
func (p *Peer) HandleAddr(payload []byte) error {
	peers, err := parseAddr(payload)
	if err != nil {
		return err
	}
	log.Debug("Addr: Number of peers received: ", len(peers))
	return nil
}

// HandleVersion
func (p *Peer) HandleVersion(payload []byte) error {
	buf, err := parser.New(payload)
	if err != nil {
		return err
	}
	version := buf.ReadUint32()
	services := buf.ReadUint64()
	ts := buf.ReadUint64()
	log.Info("Peer version: ", version, services, ts)
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
	log.Warn(fmt.Sprintf("Rejected: %s", payload))
}
