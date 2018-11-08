package network

import (
  "testing"
)

func TestNetwork(t *testing.T) {
	net := Network{}
	net.New()

	peer := Peer{}
	peer.New()

	net.NetworkVersion(0)
	t.Error()
}
