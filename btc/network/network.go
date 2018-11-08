package network

import (
	"math/rand"
	"bytes"
	"encoding/binary"
	"time"
	"fmt"
	log "github.com/sirupsen/logrus"
)
	= "37.59.38.74"

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
	ip
}

type Network struct {
	networkMagic
	version
	services
	userAgent
	port
	peers		[]Peer
}
