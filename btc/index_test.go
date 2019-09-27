package btc

import (
	//"github.com/chainswatch/bclib/serial"
	"github.com/joho/godotenv"

	"fmt"
	"os"
	"testing"
)

func TestIndex(t *testing.T) {
	/*
		tests := []struct {
			name   string
			height uint32
			hash   []byte
		}{
			{
				name:   "Genesis block",
				height: 0,
				hash:   []byte("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
			},
			{
				name:   "Last block of first blockfile",
				height: 120026,
				hash:   []byte("000000000000280e8b102c2a71efb35ee004cd560234cad5b6e8bbb380b94f9d"),
			},
		}
	*/

	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			t.Fatal(err)
		}
	}

	/*
		db, err := OpenIndexDb()
		if err != nil {
			t.Fatal(err)
		}
			for _, test := range tests {
				hash, err := serial.HexToBinary(test.hash)
				if err != nil {
					t.Fatal(err)
				}
				b, err := blockIndexRecord(db, serial.ReverseHex(hash))
				if err != nil {
					t.Fatal(err)
				}
				if b.NHeight != test.height {
					t.Fatalf("%v: Height: %d != %d",
						test.name, b.NHeight, test.height)
				}
			}
	*/

	idx, err := LoadHeaderIndex()
	if err != nil {
		t.Fatal(err)
	}

	prev := make([]byte, 32)
	//for h := uint32(0); h < uint32(len(idx)); h++ {
	for h := uint32(0); h < 530000; h++ {
		v, exist := idx[h]
		if !exist {
			t.Fatalf("Index: Stopped at height %d", h)
		}
		if h != v.NHeight {
			t.Fatalf("Wrong height: %d != %d", h, v.NHeight)
		}
		if h >= 525889 && h <= 525895 {
			fmt.Printf("%d %x, %d, \n", h, v.NStatus, v.NStatus&(blockHaveData|blockHaveUndo))
		}
		/*
			if h > 0 {
				if bytes.Compare(prev, v.HashPrev) != 0 {
					t.Errorf("Height %d: Wrong block hash: %x != %x (NVersion: %x, NBits: %x, %xx)",
						v.NHeight,
						prev,
						v.HashPrev,
						v.NVersion,
						v.NBits,
						v.HashMerkleRoot)
				}
			}
		*/
		copy(prev, v.Hash)
	}
}
