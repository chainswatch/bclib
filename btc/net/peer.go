package net

import (
	"net"
	"bufio"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Open a new connection with peer
func openConnection(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	log.Info("openConnection: ", conn.RemoteAddr())
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// New initializes peer structure
func (p *Peer) newConnection(ip string, port uint16) error {
	p.ip = net.ParseIP(ip)
	p.port = port

	rw, err := openConnection(fmt.Sprintf("%s:%d", p.ip.String(), p.port))
	if err != nil {
		return err
	}
	p.rw = rw

	p.invs = make(map[[32]byte]*inv)
	p.nextInvs = make(map[[32]byte]bool)
	return nil
}
