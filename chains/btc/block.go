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

func (btc *Btc) ParseBlockFromFile(blockFile *parser.BlockFile) {
  curPos, err := blockFile.Seek(0, 1)
  if err != nil {
    log.Warn(err)
  }

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
  err = btc.parseBlockTransactionsFromFile(blockFile)
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
  blockFile.Seek(int64(btc.NDataPos - 8), 0)

  btc.ParseBlockFromFile(blockFile)
}

func (btc *Btc) getTransactionsFromFile(pg *sqlx.DB) {
  for _, tx := range(btc.Transactions) {
    _, err := db.InsertTransaction(pg, tx, btc.HashBlock)
    if err != nil {
      log.Warn("Block Height=", btc.NHeight, err)
      break
    }
    for _, vin := range(tx.Vin) {
      _, err = db.InsertInput(pg, vin, tx.Hash)
      if err != nil {
        log.Warn(err)
        break
      }
    }
    for _, vout := range(tx.Vout) {
      _, err = db.InsertOutput(pg, vout, tx.Hash)
      if err != nil {
        log.Warn(err)
        break
      }
    }
  }
}

func (btc *Btc) GetBlock(nHeight int) {
  // Get block infos from... (in particular NFile, NDataPos and NUndoPos)

  pg := db.ConnectPg()
  defer pg.Close()

  for height := 0; height < nHeight; height++ {
    res, err := db.GetHeaderFromHeight(pg, height)
    if err != nil {
      log.Warn(err)
      return
    }
    btc.BlockHeader = res
    // log.Info(fmt.Sprintf("%+v", btc))
    btc.getBlockFromFile()
    btc.getTransactionsFromFile(pg)
  }
  // btc.getTransaction() // Using index
}
