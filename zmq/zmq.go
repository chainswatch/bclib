package zmq

import (
	zmq "github.com/pebbe/zmq4"
)

// OpenSub opens a ZMQ subscription socket
func OpenSub(zmqURL string, topics ...string) func(stop bool) ([]string, error) {
	subscriber, _ := zmq.NewSocket(zmq.SUB) // TODO: Error handling
	subscriber.Connect("tcp://" + zmqURL)

	for _, topic := range topics {
		subscriber.SetSubscribe(topic)
	}

	return func(stop bool) ([]string, error) {
		if stop {
			subscriber.Close()
			return nil, nil
		}
		msg, err := subscriber.RecvMessage(0)
		return msg, err
	}
}

// OpenPub opens a ZMQ publication socket
func OpenPub(zmqURL string) func(msg []string) error {
	publisher, _ := zmq.NewSocket(zmq.PUB) // TODO: Error handling
	publisher.Bind("tcp://" + zmqURL)

	return func(msg []string) error {
		if msg == nil {
			publisher.Close()
			return nil
		}
		_, err := publisher.SendMessage(msg[0], msg[1])
		return err
	}
}
