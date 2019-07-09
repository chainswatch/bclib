package btc

import (
	"github.com/chainswatch/bclib/parser"
	"github.com/chainswatch/bclib/serial"

	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"bytes"
	"strconv"
	"testing"
)

func TestBlock(t *testing.T) {
	tests := []struct {
		name    string
		fileNum int
		height  uint32
		hash		[]byte
	}{
		{
			name:    "Block 547066",
			fileNum: 1,
			height:  547066,
			hash:   []byte("000000000000000000187868fe3ac0d0bbd36161f0d97e3332928f2ea54baf5a"),
		},
		{
			name:    "Block 251718",
			fileNum: 251718,
			height:  251718,
			hash:   []byte("000000000000001fbc5a74fb56b1cf8949bcfad8e3ae06f2af638b94f7633fbc"),
		},
		{
			name:    "Block 265458",
			fileNum: 265458,
			height:  265458,
			hash:   []byte("000000000000000b7e48f88e86ceee3e97b4df7c139f5411d14735c1b3c36791"),
		},
	}

	for _, test := range tests {
		raw, err := ioutil.ReadFile("./testdata/rawblock_" + strconv.Itoa(test.fileNum))
		if err != nil {
			t.Fatalf("%s: %s", test.name, err)
		}
		rawBlock, err := parser.New(raw)
		if err != nil {
			t.Fatalf("%s: %s", test.name, err)
		}
		b, err := DecodeBlock(rawBlock)
		if err != nil {
			t.Fatalf("%s: %s", test.name, err)
		}

		// Test block header
		thash, err := serial.HexToBinary(test.hash)
		if err != nil {
			t.Fatalf("%s: %s", test.name, err)
		}
		bhash := serial.ReverseHex(b.Hash)
		if bytes.Compare(bhash, thash) != 0 {
			t.Errorf("%s: Wrong block header: %x != %x", test.name, bhash, thash)
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
