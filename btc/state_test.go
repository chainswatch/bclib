package btc

import "testing"

func tHeight(t *testing.T, h uint32) {
	var s uint64

	s = GetTotalSupply(h)
	tt := GetHeight(s)
	if tt != h {
		t.Errorf("Height error: %d != %d (s=%d)", tt, h, s)
	}
}

func tSupply(t *testing.T, s uint64) {
	var h uint32

	h = GetHeight(s)
	tt := GetTotalSupply(h)
	if tt != s {
		t.Errorf("Supply error: %d != %d (h=%d)", tt, s, h)
	}
}

func TestSupply(t *testing.T) {
	tHeight(t, 0)
	tHeight(t, 1)
	tHeight(t, 10)
	tHeight(t, 100)
	tHeight(t, 100000)
	tHeight(t, 200000)
	tHeight(t, 300000)
	tHeight(t, 400000)
	tHeight(t, 500000)

	tSupply(t, 50E8)
	tSupply(t, 100E8)
	tSupply(t, 10000E8)
	tSupply(t, 100000E8)
	tSupply(t, 1000000E8)
	tSupply(t, 10000000E8)
	tSupply(t, 15000000E8)
	tSupply(t, 18000000E8)
	tSupply(t, 19000000E8)
	tSupply(t, 20000000E8)
}
