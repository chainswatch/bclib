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

	net.AddPeer(peer)
	pl, err := net.NetworkVersion(0)
	if err != nil {
		t.Fatal(err)
	}

	msg, _ := net.NetworkMsg("version", pl)

	rw, err := Open(peer.ip.String() + ":" + peer.port)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rw.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	err = rw.Flush()
	if err != nil {
		t.Fatal(err)
	}

	response, err := rw.Peek(24)
	if err != nil {
		t.Fatal(err)
	}

	log.Info(fmt.Sprintf("Response: %s", response))
	t.Error()
}
