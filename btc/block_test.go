package btc

import (
	"git.posc.in/cw/bclib/serial"
	"git.posc.in/cw/bclib/parser"

	"testing"
	"io/ioutil"
	"strconv"
	"fmt"
	"os"
	"bufio"
)

func TestBlock(t *testing.T) {
	tests := []struct {
		name    string
		fileNum int
		height  uint32
	}{
		{
			name:       "Random block 1",
			fileNum:    1,
			height:     547066,
		},
		{
			name:       "Block 251718",
			fileNum:    251718,
			height:     251718,
		},
	}

	for _, test := range tests{
		raw, err := ioutil.ReadFile("./testdata/rawblock_" + strconv.Itoa(test.fileNum))
		if err != nil {
			t.Error(err)
		}
		rawBlock, err := parser.New(raw)
		if err != nil {
			t.Fatal(err)
		}
		b, err := DecodeBlock(rawBlock)
		if err != nil {
			t.Fatal(err)
		}
		if b.NTx != uint32(len(b.Txs)) {
			t.Error("Wrong number of transactions: ", b.NTx, " != ", len(b.Txs))
		}
		if b.NHeight != test.height {
			t.Error("Wrong block height: ", b.NHeight, " != ", test.height)
		}

		// Test each transaction hash
		file, err := os.Open(fmt.Sprintf("./testdata/hash_transactions_%d", test.fileNum))
		if err != nil {
			t.Error(err)
		}
		scanner := bufio.NewScanner(file)

		for i := uint32(0); i < b.NTx; i++ {
			tx := b.Txs[i]
			scanner.Scan()
			if fmt.Sprintf("%x", serial.ReverseHex(tx.Hash)) != scanner.Text() {
				t.Fatalf("%s %x != %s", test.name, serial.ReverseHex(tx.Hash), scanner.Text())
			}
		}
		file.Close()
	}
}
