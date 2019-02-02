package net

import (
	"bytes"
	"bufio"
	"time"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

func (p *Peer) waitMsg() (*Message, error) {
	data := make([]byte, 0)
	var flag bool
	for {
		r, err := p.rw.ReadBytes(byte(0xD9)) // TODO: Timeout
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
	return parseMsg(data), nil
}

func (p *Peer) handleConnection(fn apply, argFn interface{}) {
	for {
		msg, err := p.waitMsg()
		if err != nil {
			log.Warn(err)
			break
		}
		if err = fn(p, msg, argFn); err != nil {
			log.Warn(err)
			break
		}
	}
}

// Open a new connection with peer
func openConnection(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return nil, err
	}
	log.Info("openConnection: ", conn.RemoteAddr())
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

// AddPeer adds a new peer
func (n *Network) AddPeer(ip string, port uint16) error {
	peer := Peer{}
	if err := peer.newConnection(ip, port); err != nil {
		return err
	}
	n.peers = append(n.peers, peer)
	n.nPeers++
	return peer.handshake(n.version, n.services, n.userAgent)
}
