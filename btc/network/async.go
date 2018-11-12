package network

import (
	"git.posc.in/cw/watchers/serial"

	"net"
	"bufio"
	"bytes"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/tidwall/evio"
)

func EvioOpen(addr string) error {
	var events evio.Events
	events.Serving = func(srv evio.Server) (action evio.Action) {
		log.Info("Evio: Server started")
		return
	}
	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		c.SetContext(c)
		log.Info("Evio: Connection opened")
		//atomic.AddInt32(&connected, 1)
		out = []byte("sweetness\r\n")
		// opts.TCPKeepAlive = time.Minute * 5
		if c.LocalAddr() == nil {
			panic("nil local addr")
		}
		if c.RemoteAddr() == nil {
			panic("nil local addr")
		}
		return
	}
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		log.Info(in)
		out = in
		return
	}
	log.Info("Evio: tcp://" + addr)
 	err := evio.Serve(events, "tcp://" + addr + "?reuseport=true")
	log.Info("Evio: tcp://" + addr)
	log.Info(err)
	return err
}

func Open(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	log.Info("Open: ", conn.RemoteAddr())
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func (p *Peer) waitMsg() ([]byte, error) {
	msg := make([]byte, 0)
	for {
		// TODO: Timeout
		r, err := p.rw.ReadBytes(byte(0xD9))
		if err != nil {
			return nil, err
		}
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
func (n *Network) sendMsg(pid uint32, cmd string, pl []byte) ([]byte, error) {
	var sbuf [24]byte

	binary.LittleEndian.PutUint32(sbuf[0:4], n.networkMagic)
	copy(sbuf[4:16], cmd) // version
	binary.LittleEndian.PutUint32(sbuf[16:20], uint32(len(pl)))

	chksum := serial.DoubleSha256(pl[:])
	copy(sbuf[20:24], chksum[:4])

	msg := append(sbuf[:], pl...)

	p := n.peers[pid]
	log.Info(fmt.Sprintf("Sending %x", msg))
	_, err := p.rw.Write(msg)
	if err != nil {
		return nil, err
	}
	err = p.rw.Flush()
	if err != nil {
		return nil, err
	}
	response, err := p.waitMsg()
	if err != nil {
		return nil, err
	}
	return response, nil
}
