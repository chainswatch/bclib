package btc

import (
  "app/parser"
  // "github.com/jmoiron/sqlx"
  log "github.com/sirupsen/logrus"
  "time"
)

const (
  //! Magic numbers to identify start of block
  BLOCK_MAGIC_ID_BITCOIN MagicID = 0xd9b4bef9
  BLOCK_MAGIC_ID_TESTNET MagicID = 0x0709110b
)

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func (btc *Btc) ParseBlockHeaderFromFile(blockFile *parser.BlockFile) {
  btc.Length = blockFile.ReadUint32()
  btc.NVersion = blockFile.ReadInt32() // TODO: Uint32? (Always 1)
  btc.HashPrevBlock = blockFile.ReadBytes(32)
  btc.HashMerkleRoot = blockFile.ReadBytes(32)
  btc.NTime = time.Unix(int64(blockFile.ReadUint32()), 0)
  btc.TargetDifficulty = blockFile.ReadUint32() // TODO: Parse this as mantissa?
  btc.NNonce = blockFile.ReadUint32()
  btc.NTx = uint32(blockFile.ReadVarint())
}

func (btc *Btc) ParseBlockFromFile(blockFile *parser.BlockFile) {
  // Read header fields
  btc.ParseBlockHeaderFromFile(blockFile)

  // Parse transactions
  err := btc.parseBlockTransactionsFromFile(blockFile)
  if err != nil {
    // blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Warn(err)
    // return
  }
}

func (btc *Btc) getBlockFromFile(blockFile *parser.BlockFile) bool {
  blockFile.Seek(int64(btc.NDataPos), 0)

  magicID := MagicID(blockFile.ReadUint32())
  if magicID == 0 {
    log.Info("End of File")
    return true
  } else if magicID != BLOCK_MAGIC_ID_BITCOIN {
    // blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Fatal("Invalid block header: Can't find Magic ID ", magicID)
  }

  btc.ParseBlockFromFile(blockFile)
  // tmp, _ := blockFile.Seek(0, 1)
  // log.Debug("Starts at:", btc.NDataPos, ". Stops at:", tmp)
  return false
}
