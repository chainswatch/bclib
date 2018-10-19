package btc

import (
  "git.posc.in/cw/watchers/parser"
  "time"
)

// Parse raw data bytes
// https://github.com/bitcoin/bitcoin/blob/v0.15.1/src/chain.h#L387L407
func (btc *btc) parseBlockHeaderData(data []byte) {
  // Parse the raw bytes
  dataBuf := parser.NewDataBuf(data)

  // Discard first varint
  // FIXME: Not exactly sure why need to, but if we don't do this we won't get correct values
  dataBuf.ShiftVarint()

  btc.NHeight = uint32(dataBuf.ShiftVarint())
  btc.NStatus = uint32(dataBuf.ShiftVarint())
  btc.NTx = uint32(dataBuf.ShiftVarint())
  if btc.NStatus & (blockHaveData|blockHaveUndo) > 0 {
    btc.NFile = uint32(dataBuf.ShiftVarint())
  }
  if btc.NStatus & blockHaveData > 0 {
    btc.NDataPos = uint32(dataBuf.ShiftVarint())
  }
  if btc.NStatus & blockHaveUndo > 0 {
    btc.NUndoPos = uint32(dataBuf.ShiftVarint())
  }

  btc.NVersion = dataBuf.Shift32bit()
  btc.HashPrevBlock = dataBuf.ShiftBytes(32)
  btc.HashMerkleRoot = dataBuf.ShiftBytes(32)
  btc.NTime = time.Unix(int64(dataBuf.ShiftU32bit()), 0)
  btc.NBits = dataBuf.ShiftU32bit()
  btc.NNonce = dataBuf.ShiftU32bit()
}
