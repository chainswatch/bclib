package network

import (
  "testing"

	"fmt"
	log "github.com/sirupsen/logrus"
)

func TestNetwork(t *testing.T) {
	net := Network{}
	net.New()

	peer := Peer{}
	peer.New()
	rw, err := Open(peer.ip.String() + ":" + peer.port)
	if err != nil {
		t.Fatal(err)
	}
	peer.rw = rw

	net.AddPeer(peer)


	response, err := net.NetworkVersion(0)
	if err != nil {
		t.Fatal(err)
	}
	log.Info(fmt.Sprintf("response %s", response))
	log.Info(fmt.Sprintf("response %x", response))

	/*
	response, err := rw.Peek(100)
	if err != nil {
		t.Fatal(err)
	}

	log.Info(fmt.Sprintf("Response: %s", response))
	log.Info(fmt.Sprintf("Response: %x", response))
	*/
	t.Error()
}
