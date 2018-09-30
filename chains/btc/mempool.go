package btc

import (
  zmq "github.com/pebbe/zmq4"
  log "github.com/sirupsen/logrus"
  "fmt"
)

func (btc *Btc) MempoolWatcher() {
  //  Prepare our subscriber
  subscriber, _ := zmq.NewSocket(zmq.SUB)
  defer subscriber.Close()
  subscriber.Connect("tcp://91.121.87.21:28332")
  subscriber.SetSubscribe("hashtx")
  subscriber.SetSubscribe("rawtx")

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
    case "rawtx":
      log.Info(fmt.Sprintf("%s: %x", topic, body))
    default:
      log.Info("Unknown topic:", topic)
    }
  }
}
