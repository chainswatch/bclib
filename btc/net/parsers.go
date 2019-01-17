package net

import (
	"github.com/chainswatch/bclib/parser"

	"net"
	"fmt"
	"bytes"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

// ParseAddr returns peer slice from payload
// TODO: Return an Addr struct?
func parseAddr(payload []byte) ([]Peer, error) {
	buf, err := parser.New(payload)
	if err != nil {
		return nil, err
	}
	count := buf.ReadVarint()
	peers := make([]Peer, count)
	for _, peer := range peers {
		peer.timestamp = buf.ReadUint32()
		peer.services = buf.ReadBytes(8)
		peer.ip = net.IP(buf.ReadBytes(16))
		log.Info(peer.ip)
		peer.port = buf.ReadUint16()
	}
	return peers, nil
}

// ParseInv returns a slice of slices of inventories
func parseInv(payload []byte) ([][]byte, uint64, error) {
	buf, err := parser.New(payload)
	if err != nil {
		return nil, 0, err
	}
	count := buf.ReadVarint()
	inventory := make([][]byte, count)
	for i := uint64(0); i < count; i++ {
		inventory[i] = buf.ReadBytes(36) // type (4) + hash (32)
	}
	return inventory, count, nil
}

func parseMsg(data []byte) *Message {
	message := Message{}
	message.cmd = fmt.Sprintf("%s", bytes.Trim(data[:12], "\x00"))
	message.length = binary.LittleEndian.Uint32(data[12:16])
	// 16-20 = checksum
	message.payload = data[20:len(data)]
	if int(message.length) != len(message.payload) {
		log.Warn(fmt.Sprintf("parseMsg: length and len(payload) differ: %d != %d", message.length, len(message.payload)))
	}
	return &message
}

// waitMsg waits next message from peer
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
	message := parseMsg(data)
	return message, nil
}
