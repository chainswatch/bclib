package network

import (
	"net"
	"bufio"

	log "github.com/sirupsen/logrus"
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
