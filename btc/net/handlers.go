package net

import (
	"git.posc.in/cw/bclib/serial"

	"fmt"

	log "github.com/sirupsen/logrus"
)

// handleObject tx and block
func (p *Peer) handleObject(object string, payload []byte) (*inv, error) {
	var hash [32]byte
	copy(hash[:], serial.DoubleSha256(payload))
	if _, asked := p.nextInvs[hash]; !asked {
		return nil, fmt.Errorf("handleTx: Hash not found: %x", hash)
	}
	if _, exist := p.invs[hash]; exist {
		return nil, fmt.Errorf("handleTx: %s already exists", object)
	}
	invObj := &inv{object: object,
		raw: payload,
		timestamp: 0,
		fromIP: nil,
	}
	p.invs[hash] = invObj
	return invObj, nil
}

//sendheaders
//sendcmpct

//ping
func (p *Peer) handlePing(nonce []byte) {
	p.SendPong(nonce)
}

//feefilter
//addr (version >= 31402)
func (p *Peer) handleAddr(payload []byte) error {
	peers, err := ParseAddr(payload)
	if err != nil {
		return err
	}
	log.Info("Addr: ", len(peers))
	return nil
}

//inv
func (p *Peer) handleInv(payload []byte) error {
	inventory, count, err := ParseInv(payload)
	if err != nil {
		return err
	}
	p.SendGetdata(inventory, count)
	return nil
}

//reject
func (p *Peer) handleReject(payload []byte) {
	log.Info(fmt.Sprintf("%s", payload))
}
