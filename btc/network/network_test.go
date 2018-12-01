package network

import (
  "testing"

	//log "github.com/sirupsen/logrus"
)

func TestNetwork(t *testing.T) {
	/*
	tests := []msg {
		cmd:			"tx",
		length:		225,
		payload:	[]byte("01000000011d49446502dc107340eeeae9178f69d110b3fdfe9039239d5a34f222304bf9e6000000006a4730440220305aea2186628ec719a02451e170d5075dc3ad2da1fa5d9cb955b410a7637878022022c97a8d90df2ac74b764703bb8c7d5880eaaff25b9bfc173ff5be1241b6cda50121036fad7846f2cf0c76e1a75c2eb8c199bf0a6b97d32b677f593c642e886b2eb074ffffffff02d0121300000000001976a914b5d02aa30212786eff0a66391720a21678a72e9688acbfe40900000000001976a9141d021c472071e893285d2a9548754aacdf102b1f88ac00000000")
	}
	*/

	net := Network{}
	net.New()
	if err := net.AddPeer("37.59.38.74", 8333); err != nil{
		t.Fatal(err)
	}
}
