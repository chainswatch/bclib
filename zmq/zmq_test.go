package zmq

import (
	"testing"
	"time"
)

func TestZMQ(t *testing.T) {
	tests := []struct {
		msg []string
	}{
		{
			msg: []string{"hashtx", "19RteLGkZ9PPjDovVz4MJyySaeWi5QhzUL"},
		},
		{
			msg: []string{"rawtx", "14FCsRiFTHuraupmhbrnLq2y2W7JtAHJTc"},
		},
		{
			msg: []string{"hashblock", "1Fbqx2BqxNPcD6Gbnpd6XAGXpWSZ3P8xxq"},
		},
		{
			msg: []string{"rawblock", "3ACDAoJhzHxLEcDgf5vxH2VeJftNY6NANY"},
		},
	}
	writeMsg := OpenPub("*:5555")
	listenMsg := OpenSub("127.0.0.1:5555", "hashtx", "rawtx", "hashblock", "rawblock")

	go func() {
		time.Sleep(time.Second)
		for _, test := range tests {
			writeMsg(test.msg)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for _, test := range tests {
		res, err := listenMsg(false)
		if err != nil {
			t.Error(err)
		}
		if res[0] != test.msg[0] {
			t.Error("Wrong topic: ", res[0], " != ", test.msg[0])
		}
		if res[1] != test.msg[1] {
			t.Error("Wrong msg: ", res[1], " != ", test.msg[1])
		}
	}

	writeMsg(nil)
	listenMsg(true)
}
