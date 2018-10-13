package btc

import (
  "app/misc"
  "app/chains/parser"
  zmq "github.com/pebbe/zmq4"
  log "github.com/sirupsen/logrus"
  "fmt"
  "bytes"
)

func (btc *Btc) MempoolWatcher() {
  //  Prepare our subscriber
  subscriber, _ := zmq.NewSocket(zmq.SUB)
  defer subscriber.Close()
  subscriber.Connect("tcp://37.59.38.74:28332")
  subscriber.SetSubscribe("hashtx")
  subscriber.SetSubscribe("rawtx")

  rawTx := &parser.RawTx{}
  var prevHashTx []byte
  for {
    msg, err := subscriber.RecvMessage(0)
    if err != nil {
      log.Warn(err)
      break
    }
    topic := msg[0]
    body := []byte(msg[1])
    switch topic {
    case "hashtx":
      log.Info(fmt.Sprintf("%s: %x", topic, body))
      prevHashTx = body
    case "rawtx":
      log.Info(fmt.Sprintf("%s: %x", topic, body))
      rawTx = &parser.RawTx{Body: body, Pos: 0}
      tx, _ := parseTransaction(rawTx)
      putTransactionHash(tx)
      if !bytes.Equal(misc.ReverseHex(tx.Hash), prevHashTx) {
        log.Fatal("ERROR: ", misc.ReverseHex(tx.Hash), prevHashTx)
      }
    default:
      log.Info("Unknown topic:", topic)
    }
  }
}
