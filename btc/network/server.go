package network

import (
	"git.posc.in/cw/bclib/serial"

	"encoding/binary"
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// SendRawMsg sends command and payload
func (p *Peer) sendMsg(cmd string, pl []byte) error {
	var sbuf [24]byte

	binary.LittleEndian.PutUint32(sbuf[0:4], networkMagic)
	copy(sbuf[4:16], cmd) // version
	binary.LittleEndian.PutUint32(sbuf[16:20], uint32(len(pl)))

	chksum := serial.DoubleSha256(pl[:])
	copy(sbuf[20:24], chksum[:4])

	msg := append(sbuf[:], pl...)

	log.Info(fmt.Sprintf("Sending [%x] %x", sbuf, pl))
	_, err := p.rw.Write(msg)
	if err != nil {
		return err
	}
	err = p.rw.Flush()
	if err != nil {
		return err
	}
	return nil
}


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
	response, err := n.sendVersion(0)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("response %x", response))

	response, err = n.sendVerack(0)
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

	for i := 0; i < 10; i++ {
		message, err := p.waitMsg()
		if err != nil {
			log.Warn(err)
			break
		}
		log.Info(fmt.Sprintf("Received: %s %d %x", message.cmd, message.length, message.payload))
		switch message.cmd {
		case "addr":
			p.handleAddr(message.payload)
		case "ping":
			p.handlePing(message.payload)
		case "inv":
			p.handleInv(message.payload)
		case "tx":
			p.handleTx(message.payload)
		case "reject":
			p.handleReject(message.payload)
		}
	}
	return nil
}
