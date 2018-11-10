package network

import (
	"net"
	"bufio"

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

func getMsg(rw *bufio.ReadWriter) ([]byte, error) {
	msg := make([]byte, 0)
	for {
		r, err := rw.ReadBytes(byte(0xD9))
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
