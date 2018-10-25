package btc

import (
  "git.posc.in/cw/watchers/serial"
  "git.posc.in/cw/watchers/parser"

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
  }{
    {
      name:       "Random block 1",
      fileNum:    1,
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
    rawBlock := parser.New(raw)
    if rawBlock == nil {
      t.Error("parse.New(): Wrong input type. Must be []byte or uint32")
    }
    btc, _ := DecodeBlock(rawBlock)
    if btc.NTx != uint32(len(btc.Transactions)) {
      t.Error("Wrong number of transactions: ", btc.NTx, " != ", len(btc.Transactions))
    }

    // Test each transaction hash
    file, err := os.Open("./testdata/hash_transactions_1")
    if err != nil {
      t.Error(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)

    for t := uint32(0); t < btc.NTx; t++ {
      tx := btc.Transactions[t]
      scanner.Scan()
      if fmt.Sprintf("%x", serial.ReverseHex(tx.Hash)) != scanner.Text() {
        log.Info(fmt.Sprintf("%x", serial.ReverseHex(tx.Hash)), " != ", scanner.Text())
        log.Info(fmt.Sprintf("%x", tx.Vin[0].Script[0]))
      }
    }
    // t.Error()
  }
}
