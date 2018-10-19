package btc

import (
  "app/parser"
  log "github.com/sirupsen/logrus"
  "time"
)

const (
  //! Magic numbers to identify start of block

  blockMagicIDBitcoin magicID = 0xd9b4bef9
  blockMagicIDTestnet magicID = 0x0709110b
)

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func (btc *btc) parseBlockHeaderFromFile(blockFile *parser.BlockFile) {
  btc.Length = blockFile.ReadUint32()
  btc.NVersion = blockFile.ReadInt32() // TODO: Uint32? (Always 1)
  btc.HashPrevBlock = blockFile.ReadBytes(32)
  btc.HashMerkleRoot = blockFile.ReadBytes(32)
  btc.NTime = time.Unix(int64(blockFile.ReadUint32()), 0)
  btc.TargetDifficulty = blockFile.ReadUint32() // TODO: Parse this as mantissa?
  btc.NNonce = blockFile.ReadUint32()
  btc.NTx = uint32(blockFile.ReadVarint())
}

func (btc *btc) parseBlockFromFile(blockFile *parser.BlockFile) {
  // Read header fields
  btc.parseBlockHeaderFromFile(blockFile)

  // Parse transactions
  err := btc.parseBlockTransactionsFromFile(blockFile)
  if err != nil {
    // blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Warn(err)
    // return
  }
}

func (btc *btc) getBlockFromFile(blockFile *parser.BlockFile) bool {
  blockFile.Seek(int64(btc.NDataPos), 0)

  mID := magicID(blockFile.ReadUint32())
  if mID == 0 {
    log.Info("End of File")
    return true
  } else if mID != blockMagicIDBitcoin {
    // blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Fatal("Invalid block header: Can't find Magic ID ", mID)
  }

  btc.ParseBlockFromFile(blockFile)
  // tmp, _ := blockFile.Seek(0, 1)
  // log.Debug("Starts at:", btc.NDataPos, ". Stops at:", tmp)
  return false
}
