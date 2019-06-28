package parser

import (
	"github.com/chainswatch/bclib/serial"
	"testing"
)

func TestParser(t *testing.T) {
	for i := uint64(0); i < 10000; i++ {
		b := CompactSize(i)
		raw, err := New(b)
		if err != nil {
			t.Fatal(err)
		}
		j := raw.ReadCompactSize()
		if i != j {
			t.Fatalf("%d != %d", i, j)
		}
	}

	raw, err := serial.HexToBinary([]byte("fde814"))
	if err != nil {
		t.Fatal(err)
	}
	br, err := New(raw)
	if err != nil {
		t.Fatal(err)
	}
	if br.ReadCompactSize() != 5352 {
		t.Errorf("!= 5352")
	}
}
