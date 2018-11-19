package network

import (
	"git.posc.in/cw/bclib/serial"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func checkType(received string, expected string) error {
	if received != expected {
		return fmt.Errorf("checkType: Unexpected response from peer. Received %s != %s", received, expected)
	}
	return nil
}

//sendheaders
//sendcmpct
//ping
func (p *Peer) handlePing(nonce []byte) {
	p.sendPong(nonce)
}

//feefilter
//addr (version >= 31402)
func (p *Peer) handleAddr(payload []byte) {
	peers := parseAddr(payload)
	log.Info("Addr: ", len(peers))
}

//inv
func (p *Peer) handleInv(payload []byte) {
	inventory, count := parseInv(payload)
	p.sendGetdata(inventory, count)
}

// handleObject tx and block
func (p *Peer) handleObject(object string, payload []byte) {
	var hash [32]byte
	copy(hash[:], serial.DoubleSha256(payload))
	if _, asked := p.nextInvs[hash]; asked {
		if _, exist := p.invs[hash]; !exist {
			p.invs[hash] = inv{object: object,
				raw: payload,
				timestamp: 0,
				fromIP: nil}
		} else {
			log.Warn("handleTx: ", object, " already exists")
		}
	} else {
		log.Warn("handleTx: Hash not found")
	}
	// TODO: Send through ZMQ
}

// version

// verack

//reject
func (p *Peer) handleReject(payload []byte) {
	log.Info(fmt.Sprintf("%s", payload))
}
