package network

import (
	"git.posc.in/cw/watchers/parser"

	"time"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

// TODO: Return an Addr struct?
func parseAddr(payload []byte) ([][]byte, uint64) {
	buf := parser.New(payload)
	count := buf.ShiftVarint()
	res := make([][]byte, count)
	for i := uint64(0); i < count; i++ {
		timestamp := buf.ReadUint32()
		services := buf.ReadBytes(8)
		ip := net.IP(buf.ReadBytes(16))
		port := buf.ReadUint16()
	}
	return res, count
}

func parseInv(payload []byte) ([][]byte, count) {
	buf := parser.New(payload)
	count := buf.ShiftVarint()
	inventory := make([][]byte, count)
	for i := uint64(0); i < count; i++ {
		inventory[i] = buf.ReadBytes(36) // type (4) + hash (32)
		// log.Info(fmt.Sprintf("%x %x", invType, hash))
	}
}
