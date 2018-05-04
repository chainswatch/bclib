package chains

import (
  db "app/chains/repository"
  log "github.com/sirupsen/logrus"
  "fmt"
  "time"
  "encoding/hex"
)

// Parse raw data bytes
// https://github.com/bitcoin/bitcoin/blob/v0.15.1/src/chain.h#L387L407
func (btc *BtcBlock) getBlockHeader() {
  //fmt.Printf("Begin: blockHash: %v, %d bytes\n", btc.HashPrevBlock, len(btc.HashPrevBlock))

  // Get data
  data, err := btc.IndexDb.Get(append([]byte("b"), btc.HashPrevBlock...), nil)
  if err != nil {
    fmt.Println(err)
  }
  // fmt.Printf("rawBlockHeader: %v\n", data)

  // Parse the raw bytes
  dataBuf := NewDataBuf(data)
  //fmt.Printf("rawData: %v\n", b)
  //dataHex := hex.EncodeToString(b)
  //fmt.Printf("rawData: %v\n", dataHex)

  // Discard first varint
  // FIXME: Not exactly sure why need to, but if we don't do this we won't get correct values
  dataBuf.ShiftVarint()

  btc.NHeight = int32(dataBuf.ShiftVarint())
  btc.NStatus = uint32(dataBuf.ShiftVarint())
  btc.NTx = uint32(dataBuf.ShiftVarint())
  if btc.NStatus & (BLOCK_HAVE_DATA|BLOCK_HAVE_UNDO) > 0 {
    btc.NFile = int32(dataBuf.ShiftVarint())
  }
  if btc.NStatus & BLOCK_HAVE_DATA > 0 {
    btc.NDataPos = uint32(dataBuf.ShiftVarint())
  }
  if btc.NStatus & BLOCK_HAVE_UNDO > 0 {
    btc.NUndoPos = uint32(dataBuf.ShiftVarint())
  }

  btc.NVersion = dataBuf.Shift32bit()
  btc.HashBlock = append([]byte(nil), btc.HashPrevBlock...)
  btc.HashPrevBlock = dataBuf.ShiftBytes(32)
  btc.HashMerkleRoot = dataBuf.ShiftBytes(32)
  btc.NTime = time.Unix(int64(dataBuf.ShiftU32bit()), 0)
  btc.NBits = dataBuf.ShiftU32bit()
  btc.NNonce = dataBuf.ShiftU32bit()
  //fmt.Printf("%+v\n", btc)
  //fmt.Printf("End: blockHash: %v, %d bytes\n", btc.HashPrevBlock, len(btc.HashPrevBlock))
}

func (btc *BtcBlock) getBlockHashInBytes(hash []byte) {
  blockHashInBytes := make([]byte, hex.DecodedLen(len(hash)))
  n, err := hex.Decode(blockHashInBytes, hash)
  if err != nil {
    fmt.Println(err)
  }
  // Reverse hex to get the LittleEndian order
  btc.HashPrevBlock = reverseHex(blockHashInBytes[:n])
}

func (btc *BtcBlock) getBlockHeaders(nBlocks int) {
  btc.IndexDb, _ = db.OpenIndexDb(btc.DataDir) // TODO: Error handling
  defer btc.IndexDb.Close()

  pg := db.ConnectPg()
  defer pg.Close()
  for i := 0; i < nBlocks; i++ {
    btc.getBlockHeader() // TODO: Errors checks
    // Copy, then insert in DB
    _, err := db.InsertHeader(pg, btc.BlockHeader)
    if err != nil {
      log.Warn(err)
    }
  }
  /*
  var id int
  err := db.Get(&id, "SELECT count(*) FROM place")
  log.Info(id)
  */
}
