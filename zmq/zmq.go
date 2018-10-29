package zmq

import (
  zmq "github.com/pebbe/zmq4"
  log "github.com/sirupsen/logrus"
  "syscall"
)

// OpenSub opens a ZMQ subscription socket
func OpenSub(zmqURL string, topics ...string) func(stop bool) ([]string, error) {
  subscriber, _ := zmq.NewSocket(zmq.SUB)
  subscriber.Connect("tcp://" + zmqURL)
  log.Info("Listening to ", zmqURL)

  for _, topic := range topics {
    subscriber.SetSubscribe(topic)
  }

  return func(stop bool) ([]string, error) {
    if stop {
      subscriber.Close()
    }
    msg, err := subscriber.RecvMessage(0)
    if err != nil {
      if zmq.AsErrno(err) == zmq.Errno(syscall.EINTR) {
        log.Info("RecvMessage: ", err, " (EINTR, handled)")
      } else {
        log.Warn("RecvMessage: ", err)
      }
    }
    return msg, err
  }
}

// OpenPub opens a ZMQ publication socket
func OpenPub(zmqURL string) func(msg []string) error {
  publisher, _ := zmq.NewSocket(zmq.PUB)
  publisher.Bind("tcp://" + zmqURL)
  log.Info("Publishing on ", zmqURL)

  return func(msg []string) error {
    if msg == nil {
      publisher.Close()
    }
    _, err := publisher.SendMessage(msg[0], msg[1])
    if err != nil {
      log.Warn("SendMessage: ", err)
    }
    return err
  }
}
