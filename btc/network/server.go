package network

import (
	"encoding/binary"
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func parseMsg(data []byte) *msg {
	message := msg{}
	message.cmd = fmt.Sprintf("%s", bytes.Trim(data[:12], "\x00"))
	message.length = binary.LittleEndian.Uint32(data[12:16])
	// 16-20 = checksum
	message.payload = data[20:len(data)-4]
	if int(message.length) != len(message.payload) {
		log.Info(message.length, "!=", len(message.payload))
	}
	return &message
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

func (p *Peer) waitMsg() (*msg, error) {
	data := make([]byte, 0)
	for {
		// TODO: Timeout
		r, err := p.rw.ReadBytes(byte(0xD9))
		if err != nil {
			return nil, err
		}
		data = append(data, r...)
		if bytes.Contains(r, []byte{0xF9, 0xBE, 0xB4, 0xD9}) {
			if len(data) == 4 && len(r) == 4 {
				data = nil
				continue
			}
			break
		}
	}
	message := parseMsg(data)
	return message, nil
}

func (n *Network) handlePeerConnect(p Peer) error {
	n.handshake()

	for i := 0; i < 20; i++ {
		message, err := p.waitMsg()
		if err != nil {
			log.Warn(err)
			break
		}
		log.Info(fmt.Sprintf("Received: %s %d %x", message.cmd, message.length, message.payload))
		switch message.cmd {
		case "addr":
			n.handleAddr(message.payload)
			i = 20
		case "ping":
			n.msgPong(message.payload)
		case "inv":
			n.handleInv(message.payload)
		}
	}
	return nil
}
