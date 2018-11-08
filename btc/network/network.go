package network

/*
func NetworkMessage() {
	// network magic
	// command
	// payload length
	// payload checksum
	// payload
}
*/

type Peer struct {
	ip						string
	port					uint32
}

type Network struct {
	networkMagic	uint32
	version				uint32
	services			uint32
	userAgent			string
	port					uint32
	peers					[]Peer
	nPeers				uint32
}
