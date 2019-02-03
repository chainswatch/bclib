package net

import (
	"bufio"
	"time"
	"fmt"
	"net"
)

// Open a new connection with peer
func openConnection(addr string) (*bufio.ReadWriter, error) {
	dialer := &net.Dialer{
		Timeout:   1 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// newConnection initializes peer structure
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
