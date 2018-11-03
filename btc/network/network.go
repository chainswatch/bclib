package main

import (
)

const port = 8333
const peer = "37.59.38.74"

func (c *OneConnection) SendRawMsg(cmd string, pl []byte) (e error) {
	c.Mutex.Lock()

	/*if c.X.Debug {
		fmt.Println(c.ConnID, "sent", cmd, len(pl))
	}*/

	if !c.broken {
		// we never allow the buffer to be totally full because then producer would be equal consumer
		if bytes_left := SendBufSize - c.BytesToSent(); bytes_left <= len(pl) + 24 {
			c.Mutex.Unlock()
			println(c.PeerAddr.Ip(), c.Node.Version, c.Node.Agent, "Peer Send Buffer Overflow @",
				cmd, bytes_left, len(pl)+24, c.SendBufProd, c.SendBufCons, c.BytesToSent())
			c.Disconnect("SendBufferOverflow")
			common.CountSafe("PeerSendOverflow")
			return errors.New("Send buffer overflow")
		}

		c.counters["sent_"+cmd]++
		c.counters["sbts_"+cmd] += uint64(len(pl))

		common.CountSafe("sent_"+cmd)
		common.CountSafeAdd("sbts_"+cmd, uint64(len(pl)))
		var sbuf [24]byte

		c.X.LastCmdSent = cmd
		c.X.LastBtsSent = uint32(len(pl))

		binary.LittleEndian.PutUint32(sbuf[0:4], common.Version)
		copy(sbuf[0:4], common.Magic[:])
		copy(sbuf[4:16], cmd)
		binary.LittleEndian.PutUint32(sbuf[16:20], uint32(len(pl)))

		sh := btc.Sha2Sum(pl[:])
		copy(sbuf[20:24], sh[:4])

		c.append_to_send_buffer(sbuf[:])
		c.append_to_send_buffer(pl)

		if x:=c.BytesToSent(); x>c.X.MaxSentBufSize {
			c.X.MaxSentBufSize = x
		}
	}
	c.Mutex.Unlock()
	select {
		case c.writing_thread_push <- true:
		default:
	}
	return
}

func (c *OneConnection) SendVersion() {
	b := bytes.NewBuffer([]byte{})

	binary.Write(b, binary.LittleEndian, uint32(common.Version)) // Protocol version, 4 bytes, LE, 70015
	binary.Write(b, binary.LittleEndian, uint64(common.Services)) // Network services of sender
	binary.Write(b, binary.LittleEndian, uint64(time.Now().Unix())) // Timestamp

	b.Write(c.PeerAddr.NetAddr.Bytes())
	if ExternalAddrLen()>0 {
		b.Write(BestExternalAddr())
	} else {
		b.Write(bytes.Repeat([]byte{0}, 26))
	}

	b.Write(nonce[:])

	common.LockCfg()
	btc.WriteVlen(b, uint64(len(common.UserAgent)))
	b.Write([]byte(common.UserAgent))
	common.UnlockCfg()

	binary.Write(b, binary.LittleEndian, uint32(common.Last.BlockHeight())) // BlockHeight
	b.WriteByte(1)  // don't notify me about txs

	c.SendRawMsg("version", b.Bytes())
}

func main() {

}
