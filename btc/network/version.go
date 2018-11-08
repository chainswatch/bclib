package network

func NetworkVersion() {
	b := bytes.NewBuffer([]byte{})

	binary.Write(b, binary.LittleEndian, uint32(Version)) // Protocol version, 70015
	binary.Write(b, binary.LittleEndian, uint64(Services)) // Services
	binary.Write(b, binary.LittleEndian, uint64(time.Now().Unix())) // Timestamp

	// Network address of receiver (26)
	b.Write([]byte(Peer)) // Network address of receiver
	binary.Write(b, binary.LittleEndian, uint16(Port)) // Network port of receiver

	// Network address of emitter (26)
	b.Write(bytes.Repeat([]byte{0}, 26))

	binary.Write(b, binary.LittleEndian, uint64(rand.Intn(2^64))) // nonce, 8 bytes

	binary.Write(b, binary.LittleEndian, uint64(len(UserAgent)))
	b.Write([]byte(UserAgent))
	// b.Write([]byte{0})

	binary.Write(b, binary.LittleEndian, uint32(0)) // Last blockheight received
	b.WriteByte(1)  // don't notify me about txs (BIP37)

	log.Info(fmt.Sprintf("version %x", b.Bytes()))
	SendRawMsg("version", b.Bytes())
}

// SendRawMsg sends command and payload
func SendRawMsg(cmd string, pl []byte) (e error) {

	fmt.Println(c.ConnID, "sent", cmd, len(pl))

	binary.LittleEndian.PutUint32(sbuf[0:4], common.Version)
	copy(sbuf[0:4], common.Magic[:]) // Network Magic
	copy(sbuf[4:16], cmd) // version
	binary.LittleEndian.PutUint32(sbuf[16:20], uint32(len(pl)))

	sh := btc.Sha2Sum(pl[:])
	copy(sbuf[20:24], sh[:4])

	c.append_to_send_buffer(sbuf[:])
	c.append_to_send_buffer(pl) // payload

	return
}

// this function assumes that there is enough room inside sendBuf
func append_to_send_buffer(d []byte) {
 room_left := SendBufSize - c.SendBufProd
 if room_left>=len(d) {
  copy(c.sendBuf[c.SendBufProd:], d)
  room_left = c.SendBufProd+len(d)
 } else {
  copy(c.sendBuf[c.SendBufProd:], d[:room_left])
  copy(c.sendBuf[:], d[room_left:])
 }
 c.SendBufProd = (c.SendBufProd + len(d)) & SendBufMask
}
