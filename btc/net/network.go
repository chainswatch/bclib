package net

import (
	"bufio"
	"time"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
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
				break
			}
			switch m.Cmd() {
			case "addr":
				peers, err := p.handleAddr(m.Payload())
				if err != nil {
					log.Warn(err)
				}
				for _,peer := range peers {
					if err := n.AddPeer(peer); err != nil {
						log.Warn(err)
					}
				}
			case "ping":
				p.handlePing(m.Payload())
			default:
				if err = n.fn(p, m, n.argFn); err != nil {
					log.Warn(err)
					break
				}
			}
			alive <- true
			time.Sleep(10 * time.Millisecond)
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

// Open a new connection with peer
func openConnection(addr string) (*bufio.ReadWriter, error) {
	dialer := &net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// AddPeer adds a new peer
func (n *Network) AddPeer(p *Peer) error {
	ip := p.ip.String()
	if p.ip.To4() == nil {
		ip = fmt.Sprintf("[%s]", ip)
	}
	rw, err := openConnection(fmt.Sprintf("%s:%d", ip, p.port))
	if err != nil {
		return err
	}

	p.rw = rw
	p.queue = NewQueue(10000)

	if _, exists := n.peers[p.ip.String()]; exists {
		return fmt.Errorf("Already connected to that peer (%s:%d)", p.ip, p.port)
	}
	if _, exists := n.banned[p.ip.String()]; exists {
		return fmt.Errorf("Peer banned (%s:%d)", p.ip, p.port)
	}
	n.newAddr[p.ip.String()] = p
	return nil
}


// Watch network and adds new peers
func (n *Network) Watch() {
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
			n.peers[k] = p
			go n.handle(p)
		}
		time.Sleep(10 * time.Second)
	}
	log.Warn("All peers disconnected. Exiting")
}
