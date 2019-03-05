package net

import (
	"bufio"
	"fmt"
	"net"
	"time"

	zmq "github.com/pebbe/zmq4"
	log "github.com/sirupsen/logrus"
)

// SetContext sets network context
func (n *Network) SetContext(ctx *zmq.Context) {
	n.ctx = ctx
}

// SetMaxPeers sets the maximum number of peers
func (n *Network) SetMaxPeers(m uint32) {
	n.maxPeers = m
}

// ConnectedPeers returns the number of connected peers
func (n *Network) ConnectedPeers() []string {
	peers := make([]string, len(n.peers))
	i := 0
	for _, v := range n.peers {
		peers[i] = fmt.Sprintf("%s:%d", v.ip, v.port)
		i++
	}
	return peers
}

// New initializes network structure
func (n *Network) New(fn apply) {
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333

	n.peers = make(map[string]*Peer)
	n.newAddr = make(map[string]*Peer)
	n.banned = make(map[string]bool)
	n.maxPeers = 1

	n.fn = fn
}

func (n *Network) action(p *Peer, alive chan bool, kill chan bool) {
	for {
		select {
		case <-kill:
			log.Info("Child process killed properly")
			return
		default:
			m, err := p.waitMsg()
			if err != nil {
				p.errors++
				log.Warn(fmt.Sprintf("waitMsg: error no %d. %s", p.errors, err))
				if p.errors > 2 {
					log.Warn("Too many errors: disconnecting from peer")
					return
				}
				continue
			}
			if p.errors > 0 {
				p.errors--
			}
			switch m.Cmd() {
			case "addr":
				peers, err := p.handleAddr(m.Payload())
				if err != nil {
					log.Warn(err)
				}
				for _, peer := range peers {
					if err := n.AddPeer(peer); err != nil {
						log.Debug(err)
					}
				}
			case "ping":
				p.handlePing(m.Payload())
			default:
				if err = n.fn(p, m); err != nil {
					return // TODO: Relay error msg?
				}
			}
			alive <- true
		}
	}
}

// TODO: remove log
func (n *Network) handle(p *Peer) {
	alive := make(chan bool, 1)
	kill := make(chan bool, 1)
	go n.action(p, alive, kill)
	for {
		select {
		case <-alive:
		case <-time.After(3 * time.Minute):
			log.Warn("3 minutes passed without receiving any message from ", p.ip)
			n.banned[p.ip.String()] = true
			delete(n.peers, p.ip.String())
			kill <- true
			break
		}
	}
}

// AddPeer adds a new peer
func (n *Network) AddPeer(p *Peer) error {
	var err error

	if _, exists := n.peers[p.ip.String()]; exists {
		return fmt.Errorf("Already connected to that peer (%s:%d)", p.ip, p.port)
	}
	if _, exists := n.banned[p.ip.String()]; exists {
		return fmt.Errorf("Peer banned (%s:%d)", p.ip, p.port)
	}

	ip := p.ip.String()
	if p.ip.To4() == nil {
		ip = fmt.Sprintf("[%s]", ip)
	}
	dialer := &net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	p.conn, err = dialer.Dial("tcp", fmt.Sprintf("%s:%d", ip, p.port))
	if err != nil {
		return err
	}
	p.rw = bufio.NewReadWriter(bufio.NewReader(p.conn), bufio.NewWriter(p.conn))

	p.queue = NewQueue(10000)
	n.newAddr[p.ip.String()] = p
	return nil
}

// Watch network and adds new peers
func (n *Network) Watch(url string) {
	var err error
	// TODO: Use a channel
	for len(n.peers) > 0 || len(n.newAddr) > 0 {
		for k, p := range n.newAddr {
			if len(n.peers) >= int(n.maxPeers) {
				break
			}
			delete(n.newAddr, k)
			if err = p.handshake(n.version, n.services, n.userAgent); err != nil {
				log.Warn("Handshake: ", err)
				continue
			}
			if len(url) > 0 {
				log.Info("Connect ", p.ip)
				if p.Pub, err = n.ctx.NewSocket(zmq.PUB); err != nil {
					log.Warn("NewSocket: ", err)
				}
				if err = p.Pub.Connect(url); err != nil { // TODO: Disconnect properly
					log.Warn("Connect: ", err)
				}
			}
			n.peers[k] = p
			go n.handle(p)
		}
		time.Sleep(10 * time.Second)
	}
	log.Warn("All peers disconnected. Exiting")
}
