package net

// ConnectedPeers returns the number of connected peers
func (n *Network) ConnectedPeers() uint32 {
	return n.nPeers
}

// New initializes network structure
func (n *Network) New() {
	n.version = 70015
	n.services = 0
	n.userAgent = "/CW:01/"
	n.port = 8333
	n.nPeers = 0
	n.maxPeers = 10
}

// apply is passed as an argument to Watch
type apply func(*Peer, *Message, interface{}) error

// Watch connected peers and apply fn when a message is received
func (n *Network) Watch(fn apply, argFn interface{}) {
	for _,p := range n.peers {
		p := p
		go p.handle(fn, argFn)
	}
}

// AddPeer adds a new peer
func (n *Network) AddPeer(ip string, port uint16) error {
	peer := Peer{}
	if err := peer.newConnection(ip, port); err != nil {
		return err
	}
	n.peers = append(n.peers, peer)
	n.nPeers++
	return peer.handshake(n.version, n.services, n.userAgent)
}
