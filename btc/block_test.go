package btc

import (
	"git.posc.in/cw/bclib/models"
	"git.posc.in/cw/bclib/serial"
	"git.posc.in/cw/bclib/parser"

	"testing"
	"io/ioutil"
	"strconv"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	}

	for _, test := range tests{
		hash, err := ioutil.ReadFile("./testdata/hashblock_" + strconv.Itoa(test.fileNum))
		if err != nil {
			t.Error(err)
		}
		log.Info(fmt.Sprintf("Block Hash: %x", hash))
		raw, err := ioutil.ReadFile("./testdata/rawblock_" + strconv.Itoa( test.fileNum))
		if err != nil {
			t.Error(err)
		}
		rawBlock, err := parser.New(raw)
		if err != nil {
			t.Fatal(err)
		}
		b := &models.Block{}
		err = DecodeBlock(b, rawBlock)
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
		file, err := os.Open("./testdata/hash_transactions_1")
		if err != nil {
			t.Error(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		for i := uint32(0); i < b.NTx; i++ {
			tx := b.Txs[i]
			scanner.Scan()
			if fmt.Sprintf("%x", serial.ReverseHex(tx.Hash)) != scanner.Text() {
				t.Errorf("%x", tx.Vin[0].Script[0])
				t.Fatalf("%x != %s", serial.ReverseHex(tx.Hash), scanner.Text())
			}
		}
	}
}
