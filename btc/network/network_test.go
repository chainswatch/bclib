package network

import (
  "testing"

	//log "github.com/sirupsen/logrus"
)

func TestNetwork(t *testing.T) {

	net := Network{}
	net.New()

	net.NewPeer("37.59.38.74", "8333")

	t.Error()
}
