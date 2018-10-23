package btc

import (
  "git.posc.in/cw/watchers/parser"

  log "github.com/sirupsen/logrus"
  "time"
)

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func ParseBlockHeaderFromFile(btc *btc, br parser.Reader) {
  btc.Length = br.ReadUint32()
  btc.NVersion = br.ReadInt32() // TODO: Uint32? (Always 1)
  btc.HashPrevBlock = br.ReadBytes(32)
  btc.HashMerkleRoot = br.ReadBytes(32)
  btc.NTime = time.Unix(int64(br.ReadUint32()), 0)
  btc.TargetDifficulty = br.ReadUint32() // TODO: Parse this as mantissa?
  btc.NNonce = br.ReadUint32()
  btc.NTx = uint32(br.ReadVarint())
}

// Parse raw data bytes
// https://github.com/bitcoin/bitcoin/blob/v0.15.1/src/chain.h#L387L407
func decodeBlockHeader(btc *btc, br parser.Reader) {
  // Discard first varint
  // FIXME: Not exactly sure why need to, but if we don't do this we won't get correct values
  br.ReadVarint()

  btc.NHeight = uint32(br.ReadVarint())
  btc.NStatus = uint32(br.ReadVarint())
  btc.NTx = uint32(br.ReadVarint())
  if btc.NStatus & (blockHaveData|blockHaveUndo) > 0 {
    btc.NFile = uint32(br.ReadVarint())
  }
  if btc.NStatus & blockHaveData > 0 {
    btc.NDataPos = uint32(br.ReadVarint())
  }
  if btc.NStatus & blockHaveUndo > 0 {
    btc.NUndoPos = uint32(br.ReadVarint())
  }

  btc.NVersion = br.ReadInt32()
  btc.HashPrevBlock = br.ReadBytes(32)
  btc.HashMerkleRoot = br.ReadBytes(32)
  btc.NTime = time.Unix(int64(br.ReadUint32()), 0)
  btc.NBits = br.ReadUint32()
  btc.NNonce = br.ReadUint32()
}

func DecodeBlock(br parser.Reader) {
  // decodeBlockHeader(br)

  _, err := DecodeTx(br)
  if err != nil {
    // blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Warn(err)
    // return
  }
}
