package net

import (
	"git.posc.in/cw/bclib/parser"

	"net"
	"fmt"
	"bytes"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

// TODO: Return an Addr struct?
// ParseAddr returns peer slice from payload
func ParseAddr(payload []byte) ([]Peer, error) {
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
		peer.port = buf.ReadUint16()
	}
	return peers, nil
}

// ParseInv returns a slice of slices of inventories
func ParseInv(payload []byte) ([][]byte, uint64, error) {
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
	message.payload = data[20:len(data)-4]
	if int(message.length) != len(message.payload) {
		log.Info(message.length, "!=", len(message.payload))
	}
	return &message
}

// WaitMsg waits next message from peer
func (p *Peer) WaitMsg() (*Message, error) {
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
