package network

import (
	"net"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// New initializes peer structure
func (p *Peer) New(ip, port string) error {
	p.ip = net.ParseIP(ip)
	p.port = port

	rw, err := Open(p.ip.String() + ":" + p.port)
	if err != nil {
		return err
	}
	p.rw = rw
	return nil
}

func (n *Network) handshake() error {
	response, err := n.msgVersion(0)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("response %x", response))

	response, err = n.msgVerack(0)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("response %x", response))
	return nil
}

func (n *Network) handlePeerConnect(p Peer) error {
	n.handshake()

	for {
		response, err := p.waitMsg()
		if err != nil {
			log.Warn(err)
			break
		}
		log.Info(fmt.Sprintf("response %s", response))
		log.Info(fmt.Sprintf("response %x", response))
	}
	return nil
}

func (n *Network) NewPeer(ip, port string) error {
	peer := Peer{}
	if err := peer.New(ip, port); err != nil {
		return err
	}
	n.peers = append(n.peers, peer)
	n.nPeers++
	n.handlePeerConnect(peer)
	return nil
}

// New initializes network structure
func (n *Network) New() {
	n.networkMagic = 0xD9B4BEF9 // Maybe LE
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333
	n.nPeers = 0
}
