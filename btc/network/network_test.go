package network

import (
  "testing"

	"time"
	log "github.com/sirupsen/logrus"
)

func TestNetwork(t *testing.T) {

	go EvioOpen("127.0.0.1" + ":" + "8333")

	time.Sleep(time.Second)
	log.Info("====================================================")

	net := Network{}
	net.New()

	net.NewPeer("37.59.38.74", "8333")

	t.Error()
}
