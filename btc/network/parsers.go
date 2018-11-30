package network

import (
	"git.posc.in/cw/bclib/parser"

	"net"
)

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
		peer.port = buf.ReadUint16()
	}
	return peers, nil
}

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
