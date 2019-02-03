package net

import (
	"bytes"
	"fmt"

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
	log.Info(fmt.Sprintf("%s parseMsg", p.ip))
	return parseMsg(data)
}

// TODO: remove log
func (p *Peer) handle(fn apply, argFn interface{}) {
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
