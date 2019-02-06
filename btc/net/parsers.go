package net

import (
	"github.com/chainswatch/bclib/parser"

	"bytes"
	"encoding/binary"
	"fmt"
	"net"

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
		buf.ReadUint32() // time (!= timestamp)
		peer.services = buf.ReadUint64()
		peer.ip = net.IP(buf.ReadBytes(16))
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

func parseMsg(data []byte) (*Message, error) {
	message := Message{}
	log.Info("Msg Length ", len(data))
	message.cmd = fmt.Sprintf("%s", bytes.Trim(data[:12], "\x00"))
	message.length = binary.LittleEndian.Uint32(data[12:16])
	// TODO: 16-20 = checksum
	message.payload = data[20:len(data)]
	if int(message.length) != len(message.payload) {
		return nil, fmt.Errorf("parseMsg: length and len(payload) differ: %d != %d", message.length, len(message.payload))
	}
	return &message, nil
}
