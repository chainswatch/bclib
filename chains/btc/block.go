package btc

import (
  "app/chains/parser"
  db "app/chains/repository"
  "database/sql"
  // "github.com/jmoiron/sqlx"
  log "github.com/sirupsen/logrus"
  "time"
)

const (
  //! Magic numbers to identify start of block
  BLOCK_MAGIC_ID_BITCOIN MagicID = 0xd9b4bef9
  BLOCK_MAGIC_ID_TESTNET MagicID = 0x0709110b
)

func (btc *Btc) insertBlock(pgtx *sql.Tx) {
  db.InsertHeader(pgtx, btc.Block)

  nextTx := db.PrepareInsertTransaction(pgtx)
  nextVin := db.PrepareInsertInput(pgtx)
  nextVout := db.PrepareInsertOutput(pgtx)
  for _, tx := range(btc.Transactions) {
    nextTx(tx)
    for _, vin := range(tx.Vin) {
      nextVin(vin, tx.Hash)
    }
    for _, vout := range(tx.Vout) {
      nextVout(vout, tx.Hash)
    }
  }
}

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

func (btc *Btc) GetAllBlocks() {
  var height uint32
  res, err := db.GetLastBlockHeader(btc.SqlDb)
  if err != nil {
    log.Warn(err)
    btc.NDataPos = 0
    btc.NFile = 0
    height = 0
  } else {
    btc.BlockHeader = res
    height = btc.NHeight + 1
    btc.NDataPos += btc.Length + 8
    log.Info("Starting from file ", btc.NFile, " at pos ", btc.NDataPos, " and height ", height)
  }

  // Loop through files
  for {
    log.Info("Reading file ", btc.NFile)
    blockFile, err := parser.NewBlockFile(btc.DataDir, btc.NFile)
    if err != nil {
      log.Warn(err)
      break
    }

    // Open pgSQL transaction
    tx, err := btc.SqlDb.Begin()
    if err != nil {
      log.Fatal("Begin:", err)
    }

    // loop through blocks
    for {
      btc.NHeight = uint32(height)
      if btc.getBlockFromFile(blockFile) {
        break
      }
      btc.insertBlock(tx)
      if height % 10000 == 0 {
        log.Info(height)
      }
      btc.NDataPos += btc.Length + 8 // Jump to next block
      height++
    }

    // Close pgSQL transaction
    err = tx.Commit()
    if err != nil {
      log.Fatal(err)
    }

    blockFile.Close()
    btc.NFile++
    btc.NDataPos = 0
  }
  // btc.getTransaction() // Using index
}
