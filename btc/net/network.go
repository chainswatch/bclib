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
func (n *Network) New(fn apply, argFn interface{}) {
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333

	n.peers = make(map[string]*Peer)
	n.newAddr = make(map[string]*Peer)
	n.banned = make(map[string]bool)
	n.maxPeers = 10

	n.fn = fn
	n.argFn = argFn
}

// TODO: remove log
func (n *Network) handle(p *Peer) {
	chAction := make(chan bool, 1)
	go p.action(chAction, n.fn, n.argFn)
	for {
		select {
		case <-chAction:
		case <-time.After(60 * time.Second):
			log.Warn("60 seconds passed without receiving any message from ", p.ip)
			n.banned[p.ip.String()] = true
			delete(n.peers, p.ip.String())
			break
		}
	}
}

// AddPeer adds a new peer
func (n *Network) AddPeer(ip string, port uint16) error {
	p:= Peer{}
	if err := p.new(ip, port); err != nil {
		return err
	}
	n.newAddr[p.ip.String()] = &p
	return nil
}


// Watch network and adds new peers
func (n *Network) Watch() {
	// TODO: Use a channel
	for {
		log.Info(len(n.newAddr))
		for k,p := range n.newAddr {
			if len(n.peers) < int(n.maxPeers) {
				if err := p.handshake(n.version, n.services, n.userAgent); err != nil {
					log.Warn(err)
					continue
				}
				delete(n.newAddr,k)
				n.peers[k] = p
				go n.handle(p)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
