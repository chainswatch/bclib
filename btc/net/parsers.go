package net

import (
	"github.com/chainswatch/bclib/parser"

	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// ParseAddr returns peer slice from payload
// TODO: Return an Addr struct?
func parseAddr(payload []byte) ([]*Peer, error) {
	buf, err := parser.New(payload)
	if err != nil {
		return nil, err
	}
	count := buf.ReadVarint()
	peers := make([]*Peer, 0, count)
	for i := uint64(0); i < count; i++ {
		buf.ReadUint32() // time (!= timestamp)
		services := buf.ReadUint64()
		ip := net.IP(buf.ReadBytes(16))
		port := binary.BigEndian.Uint16(buf.ReadBytes(2))
		peer := NewPeer(ip.String(), port)
		peer.services = services
		peers = append(peers, peer)
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

func parseMsg(data []byte) (*Message, error) {
	message := Message{}
	message.cmd = fmt.Sprintf("%s", bytes.Trim(data[:12], "\x00"))
	message.length = binary.LittleEndian.Uint32(data[12:16])
	// TODO: 16-20 = checksum
	message.payload = data[20:]
	if int(message.length) != len(message.payload) {
		return nil, fmt.Errorf("parseMsg: length and len(payload) differ: %d != %d", message.length, len(message.payload))
	}
	return &message, nil
}
