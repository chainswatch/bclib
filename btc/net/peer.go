package net

import (
	"bytes"
	"net"
	"time"
)

// GetIP returns IP of connected peer
func (p *Peer) GetIP() string {
	return p.ip.String()
}

func (p *Peer) waitMsg() (*Message, error) {
	data := make([]byte, 0)
	var flag bool

	p.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	for {
		r, err := p.rw.ReadBytes(byte(0xD9))
		if err != nil {
			if err.Error() != "EOF" {
				return nil, err
			}
			flag = true
		}
		if bytes.Contains(r, []byte{0xF9, 0xBE, 0xB4, 0xD9}) { // TODO: Improve
			r = r[:len(r)-4]
			flag = true
		}
		if flag {
			if len(r) != 0 {
				data = append(data, r...)
			}
			if len(data) != 0 {
				break
			}
			flag = false
			continue
		}
		data = append(data, r...)
	}
	return parseMsg(data)
}

// NewPeer initializes peer structure
func NewPeer(ip string, port uint16) *Peer {
	p := &Peer{}
	p.ip = net.ParseIP(ip)
	p.port = port

	return p
}
