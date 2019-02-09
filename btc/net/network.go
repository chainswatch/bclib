package net

import (
	"time"
	"fmt"

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

func (n *Network) action(p *Peer, c chan bool) {
	for {
		m, err := p.waitMsg()
		if err != nil {
			log.Warn(err)
			break
		}
		switch m.Cmd() {
		case "addr":
			peers, err := p.handleAddr(m.Payload())
			if err != nil {
				log.Warn(err)
			}
			for _,peer := range peers {
				log.Info("New Addr: ", peer.ip.String())
				if err := n.AddPeer(peer.ip.String(), peer.port); err != nil {
					log.Warn(err)
				}
			}
		case "ping":
			p.handlePing(m.Payload())
		}
		if err = n.fn(p, m, n.argFn); err != nil {
			log.Warn(err)
			break
		}
		c <- true
	}
}

// TODO: remove log
func (n *Network) handle(p *Peer) {
	chAction := make(chan bool, 1)
	go n.action(p, chAction)
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
	if _, exists := n.peers[p.ip.String()]; exists {
		return fmt.Errorf("Already connected to that peer (%s)", ip)
	}
	if _, exists := n.banned[p.ip.String()]; exists {
		return fmt.Errorf("Peer banned (%s)", ip)
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
