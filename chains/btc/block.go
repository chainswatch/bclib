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
  pgtx, err := pg.Begin() // TODO: Error
  if err != nil {
    log.Fatal("Begin:", err)
  }
  db.InsertHeader(pgtx, btc.BlockHeader)
  for _, tx := range(btc.Transactions) {
    db.InsertTransaction(pgtx, tx, btc.NHeight)
    for _, vin := range(tx.Vin) {
      db.InsertInput(pgtx, vin, tx.Hash)
    }
    for _, vout := range(tx.Vout) {
      db.InsertOutput(pgtx, vout, tx.Hash)
    }
  }
  err = pgtx.Commit() // TODO: Error
  if err != nil {
    log.Fatal(err)
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
    log.Fatal("Invalid block header: Can't find Magic ID")
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
  } else {
    btc.BlockHeader = res
  }
  height := btc.NHeight + 1
  for ; height < maxHeight; height++ {
    btc.NHeight = uint32(height)
    btc.getBlockFromFile()
    btc.insertBlock(btc.SqlDb)
    if height % 10000 == 0 {
      log.Info(height)
    }
  }
  // btc.getTransaction() // Using index
}
