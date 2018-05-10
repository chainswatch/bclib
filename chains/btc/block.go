package btc

import (
  "app/chains/parser"
  db "app/chains/repository"
  "github.com/jmoiron/sqlx"
  log "github.com/sirupsen/logrus"
  "time"
)

const (
  //! Magic numbers to identify start of block
  BLOCK_MAGIC_ID_BITCOIN MagicID = 0xd9b4bef9
  BLOCK_MAGIC_ID_TESTNET MagicID = 0x0709110b
)

func (btc *Btc) insertBlock(pg *sqlx.DB) {
  _, err := db.InsertHeader(pg, btc.BlockHeader)
  if err != nil {
    log.Warn("Block Height=", btc.NHeight, err)
    return
  }
  for _, tx := range(btc.Transactions) {
    err := db.InsertTransaction(pg, tx, btc.NHeight)
    if err != nil {
      log.Warn("Block Height=", btc.NHeight, err)
      break
    }
    for _, vin := range(tx.Vin) {
      err = db.InsertInput(pg, vin, tx.Hash)
      if err != nil {
        log.Panic(err)
        break
      }
    }
    for _, vout := range(tx.Vout) {
      err = db.InsertOutput(pg, vout, tx.Hash)
      if err != nil {
        log.Panic(err)
        break
      }
    }
  }
}

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func (btc *Btc) ParseBlockHeaderFromFile(blockFile *parser.BlockFile) {
  btc.Length = blockFile.ReadUint32()
  btc.NVersion = blockFile.ReadInt32()
  btc.HashPrevBlock = blockFile.ReadBytes(32)
  btc.HashMerkleRoot = blockFile.ReadBytes(32)
  btc.NTime = time.Unix(int64(blockFile.ReadUint32()), 0)
  btc.TargetDifficulty = blockFile.ReadUint32() // TODO: Parse this as mantissa?
  btc.NNonce = blockFile.ReadUint32()
}

func (btc *Btc) ParseBlockFromFile(blockFile *parser.BlockFile, curPos int64) {
  blockFile.Seek(curPos, 0)
  // Read and validate Magic ID
  btc.MagicID = MagicID(blockFile.ReadUint32())
  if btc.MagicID != BLOCK_MAGIC_ID_BITCOIN {
    blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Warn("Invalid block header: Can't find Magic ID")
    return
  }

  // Read header fields
  btc.ParseBlockHeaderFromFile(blockFile)

  // Parse transactions
  err := btc.parseBlockTransactionsFromFile(blockFile)
  if err != nil {
    blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Warn(err)
    return
  }
}

func (btc *Btc) getBlockFromFile() {
  blockFile, err := parser.NewBlockFile(btc.DataDir, btc.NFile)
  if err != nil {
    log.Warn(err)
    return
  }
  defer blockFile.Close()

  // Seek to pos - 8 to start reading from block header
  // log.Info("Seeking to block at position " , btc.NDataPos)

  btc.ParseBlockFromFile(blockFile, int64(btc.NDataPos - 8))
  tmp, _ := blockFile.Seek(0, 1)
  log.Debug("Starts at:", btc.NDataPos - 8, ". Stops at:", tmp)
  btc.NDataPos = uint32(tmp + 8)
}

func (btc *Btc) GetAllBlocks(maxHeight uint32) {
  btc.NDataPos = 8
  btc.NFile = 0
  res, err := db.GetLastBlockHeader(btc.SqlDb)
  if err != nil {
    log.Warn(err)
  }
  btc.BlockHeader = res
  height := btc.NHeight + 1
  for ; height < maxHeight; height++ {
    btc.NHeight = uint32(height)
    btc.getBlockFromFile()
    btc.insertBlock(btc.SqlDb)
  }
  // btc.getTransaction() // Using index
}
