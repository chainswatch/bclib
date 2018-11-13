package network

import (
	"net"
	"bufio"

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
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333
	n.nPeers = 0
}

func Open(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	log.Info("Open: ", conn.RemoteAddr())
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}
