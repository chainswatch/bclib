package network

import (
	"git.posc.in/cw/watchers/serial"

	"net"
	"bufio"
	"bytes"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

func Open(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", addr)
	log.Info("Open: ", conn.RemoteAddr().String())
	log.Info("Open: ", conn.RemoteAddr())
	if err != nil {
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func (p *Peer) getMsg() ([]byte, error) {
	msg := make([]byte, 0)
	for {
		r, err := p.rw.ReadBytes(byte(0xD9))
		if err != nil {
			return nil, err
		}
		log.Info("getMsg:", r)
		msg = append(msg, r...)
		if bytes.Contains(r, []byte{0xF9, 0xBE, 0xB4, 0xD9}) {
			if len(msg) == 4 && len(r) == 4 {
				msg = nil
				continue
			}
			break
		}
	}
	return msg[:len(msg)-4], nil
}

// SendRawMsg sends command and payload
func (n *Network) networkMsg(pid uint32, cmd string, pl []byte) ([]byte, error) {
	var sbuf [24]byte

	binary.LittleEndian.PutUint32(sbuf[0:4], n.networkMagic)
	copy(sbuf[4:16], cmd) // version
	binary.LittleEndian.PutUint32(sbuf[16:20], uint32(len(pl)))

	chksum := serial.DoubleSha256(pl[:])
	copy(sbuf[20:24], chksum[:4])

	msg := append(sbuf[:], pl...)

	p := n.peers[pid]
	_, err := p.rw.Write(msg)
	if err != nil {
		return nil, err
	}
	err = p.rw.Flush()
	if err != nil {
		return nil, err
	}
	response, err := p.getMsg()
	if err != nil {
		return nil, err
	}
	return response, nil
}
