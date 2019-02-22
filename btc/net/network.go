package net

import (
	"bufio"
	"time"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	zmq "github.com/pebbe/zmq4"
)

// ConnectedPeers returns the number of connected peers
func (n *Network) ConnectedPeers() []string {
	peers := make([]string, len(n.peers))
	i := 0
	for _,v := range n.peers {
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
	n.maxPeers = 10

	n.fn = fn
}

func (n *Network) action(p *Peer, alive chan bool, kill chan bool) {
	for {
		select {
		case <- kill:
			log.Info("Child process killed properly")
			return
		default:
			m, err := p.waitMsg()
			if err != nil {
				log.Warn(err)
				p.errors++
				if p.errors > 10 {
					log.Warn("Too many errors: disconnecting from peer")
					return
				}
				continue
			}
			switch m.Cmd() {
			case "addr":
				peers, err := p.handleAddr(m.Payload())
				if err != nil {
					log.Warn(err)
				}
				for _,peer := range peers {
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
	// TODO: Use a channel
	for len(n.peers) > 0 || len(n.newAddr) > 0 {
		for k,p := range n.newAddr {
			if len(n.peers) >= int(n.maxPeers) {
				break
			}
			if err := p.handshake(n.version, n.services, n.userAgent); err != nil {
				log.Warn(err)
				continue
			}
			delete(n.newAddr,k)
			if len(url) > 0 {
				pub, _ := zmq.NewSocket(zmq.PUB) // TODO: Error handling
				pub.Connect(url)
				p.Pub = pub
			}
			n.peers[k] = p
			go n.handle(p)
		}
		time.Sleep(10 * time.Second)
	}
	log.Warn("All peers disconnected. Exiting")
}
