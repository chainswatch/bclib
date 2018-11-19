package network

import (
	"git.posc.in/cw/bclib/parser"

	"net"
)

// TODO: Return an Addr struct?
func parseAddr(payload []byte) []Peer {
	buf := parser.New(payload)
	count := buf.ShiftVarint()
	peers := make([]Peer, count)
	for _, peer := range peers {
		peer.timestamp = buf.ReadUint32()
		peer.services = buf.ReadBytes(8)
		peer.ip = net.IP(buf.ReadBytes(16))
		peer.port = buf.ReadUint16()
	}
	return peers
}

func parseInv(payload []byte) ([][]byte, uint64) {
	buf := parser.New(payload)
	count := buf.ShiftVarint()
	inventory := make([][]byte, count)
	for i := uint64(0); i < count; i++ {
		inventory[i] = buf.ReadBytes(36) // type (4) + hash (32)
	}
	return inventory, count
}
