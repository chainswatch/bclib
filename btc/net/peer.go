package net

import (
	"bytes"
	"bufio"
	"time"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

func (p *Peer) GetIP() string {
	return p.ip.String()
}

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
	return parseMsg(data)
}

func (p *Peer) action(c chan bool, fn apply, argFn interface{}) {
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
		c <- true
	}
}

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
func (p *Peer) new(ip string, port uint16) error {
	p.ip = net.ParseIP(ip)
	p.port = port

	rw, err := openConnection(fmt.Sprintf("%s:%d", p.ip.String(), p.port))
	if err != nil {
		return err
	}
	p.rw = rw

	p.queue = NewQueue(10000)
	return nil
}
