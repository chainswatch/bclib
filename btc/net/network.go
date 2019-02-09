package net

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// ConnectedPeers returns the number of connected peers
func (n *Network) ConnectedPeers() int {
	return len(n.peers)
}

// New initializes network structure
func (n *Network) New() {
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333

	n.peers = make(map[string]*Peer)
	n.maxPeers = 10
}

// TODO: remove log
func (n *Network) handle(p *Peer, fn apply, argFn interface{}) {
	chAction := make(chan bool, 1)
	go p.action(chAction, fn, argFn)
	for {
		select {
		case <-chAction:
		case <-time.After(60 * time.Second):
			log.Warn("60 seconds passed without receiving any message from ", p.ip)
			delete(n.peers, p.ip.String())
			break
		}
	}
}

// apply is passed as an argument to Watch
type apply func(*Peer, *Message, interface{}) error

// AddPeer adds a new peer
func (n *Network) AddPeer(ip string, port uint16, fn apply, argFn interface{}) error {
	p:= Peer{}
	if err := p.new(ip, port); err != nil {
		return err
	}
	n.peers[p.ip.String()] = &p
	if err := p.handshake(n.version, n.services, n.userAgent); err != nil {
		return err
	}
	go n.handle(&p, fn, argFn)
	return nil
}
