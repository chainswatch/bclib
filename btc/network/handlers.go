package network

import (
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

//pong

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

//
func (p *Peer) handleTx(payload []byte) {
	// Check if already exists
	// If not then send
	// ZMQ
}

// version

// verack

//reject
func (p *Peer) handleReject(payload []byte) {
	log.Info(fmt.Sprintf("%s", payload))
}
