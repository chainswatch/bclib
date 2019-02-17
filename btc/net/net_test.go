package net

import (
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

var nInv, nTx int

func handlePeers(p *Peer, m *Message) error {
	// log.Info(fmt.Sprintf("%s Received: %s %d %x", p.ip, m.Cmd(), m.Length(), m.Payload()))
	switch m.Cmd() {
	case "inv":
		nInv++
		return p.HandleInv(m.Payload())
	case "tx":
		nTx++
		_, err := p.HandleObject("tx", m.Payload())
		return err
	case "reject":
		return fmt.Errorf("Reject error")
	default:
	}
	return nil
}

func TestNetworkOne(t *testing.T) {
	n := Network{}
	n.New(handlePeers)
	if err := n.AddPeer(NewPeer("0.0.0.0", 8333)); err == nil {
		t.Fatal(err)
	}

	con := []string{"37.59.38.74","112.119.69.152","72.5.72.15", "86.97.172.251", "47.225.21.79"}
	for _, c := range con { 
		if err := n.AddPeer(NewPeer(c, 8333)); err == nil {
			break
		}
	}

	go n.Watch("")
	time.Sleep(1 * time.Second)

	if len(n.ConnectedPeers()) == 0 {
		t.Fatal("Could not connect to a single peer")
	}

	time.Sleep(20 * time.Second)

	if nInv == 0 {
		t.Error("nInv = 0")
	}
	if nTx == 0 {
		t.Error("nTx = 0")
	}
	log.Info(fmt.Sprintf("%d %d", nInv, nTx))
}

func TestNetworkMultiple(t *testing.T) {
	n := Network{}
	n.New(handlePeers)

	con := []string{"37.59.38.74","112.119.69.152","72.5.72.15", "86.97.172.251", "47.225.21.79"}
	count := 0
	for _, c := range con { 
		if err := n.AddPeer(NewPeer(c, 8333)); err == nil {
			count++
		}
		if count == 2 {
			break
		}
	}
	if count < 2 {
		t.Fatal("Failed to connect to enough peers")
	}

	time.Sleep(10 * time.Second)
	go n.Watch("")
	time.Sleep(1 * time.Second)
	if len(n.ConnectedPeers()) < 2 {
		t.Fatal("Could not connect to enough peers")
	}


	time.Sleep(30 * time.Second)

	if nInv == 0 {
		t.Error("nInv = 0")
	}
	if nTx == 0 {
		t.Error("nTx = 0")
	}
	time.Sleep(2 * time.Minute)
	log.Info(fmt.Sprintf("Connected peers: %d (%d %d)", len(n.ConnectedPeers()), nInv, nTx))
	if len(n.ConnectedPeers()) < 3 {
		t.Error("Unable to auto-connect to more peers")
	}
}

/*
func TestIPv4v6(t *testing.T) {
	n := Network{}
	n.New(handlePeers)
	n.AddPeer(NewPeer("37.59.38.74", 8333))
	n.AddPeer(NewPeer("2a0a:a545:252c:0:8cd2:921d:f532:5e42", 8333))
	t.Error()
}

func TestNetworkIPv6(t *testing.T) {
	ips := []string{"0:ffff:253b:264a:3833:3333::%lo0",
	"2a0a:a545:252c:0:8cd2:921d:f532:5e42%eth0",
	"2001:0:9d38:6abd:249e:2650:bc43:666f%wlan0",
	"2601:641:480:6340:a492:630e:bc99:7061",
	"fd87:d87e:eb43:d8cf:4827:47ac:2da6:ff74"}
	port := 8333
	for _,ip := range ips {
		if _, err := openConnection(fmt.Sprintf("[%s]:%d", ip, port)); err != nil {
			t.Error(err)
		}
	}
}

func TestNetworkZero(t *testing.T) {
	n := Network{}
	n.New(handlePeers, nil)
	go n.Watch()
	time.Sleep(5 * time.Second)
	t.Error()
}
*/
