package chains

import (
  db "app/chains/repository"
  log "github.com/sirupsen/logrus"
  "os"
  "fmt"
  "time"
)

const (
  //! Magic numbers to identify start of block
  BLOCK_MAGIC_ID_BITCOIN MagicID = 0xd9b4bef9
  BLOCK_MAGIC_ID_TESTNET MagicID = 0x0709110b
)

type BlockFile struct {
  file    *os.File
  FileNum uint32
}

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func (btc *BtcBlock) ParseBlockHeaderFromFile() {
  btc.Length = btc.ReadUint32()
  btc.NVersion = btc.ReadInt32()
  btc.HashPrevBlock = btc.ReadBytes(32)
  btc.HashMerkleRoot = btc.ReadBytes(32)
  btc.NTime = time.Unix(int64(btc.ReadUint32()), 0)
  btc.TargetDifficulty = btc.ReadUint32() // TODO: Parse this as mantissa?
  btc.NNonce = btc.ReadUint32()
}

func (btc *BtcBlock) ParseBlockFromFile() {
  curPos, err := btc.Seek(0, 1)
  if err != nil {
    log.Warn(err)
  }

  // Read and validate Magic ID
  btc.MagicID = MagicID(btc.ReadUint32())
  if btc.MagicID != BLOCK_MAGIC_ID_BITCOIN {
    btc.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Warn("Invalid block header: Can't find Magic ID")
    return
  }

  // Read header fields
  btc.ParseBlockHeaderFromFile()

  // Parse transactions
  err = btc.parseBlockTransactionsFromFile()
  if err != nil {
    btc.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Warn(err)
    return
  }
}

func (btc *BtcBlock) getBlockFile() error {
  filepath := fmt.Sprintf(btc.DataDir + "/blocks/blk%05d.dat", btc.NFile)
  log.Info("Opening file ", filepath)

  var err error
  btc.file, err = os.OpenFile(filepath, os.O_RDONLY, 0666) // TODO: Error
  return err
}

func (btc *BtcBlock) getBlockFromFile() {
  err := btc.getBlockFile()
  if err != nil {
    log.Warn(err)
    return
  }
  defer btc.file.Close()

  // Seek to pos - 8 to start reading from block header
  log.Info("Seeking to block at position " , btc.NDataPos)
  btc.Seek(int64(btc.NDataPos - 8), 0)

  btc.ParseBlockFromFile()
}

func (btc *BtcBlock) getBlock(nHeight int) {
  // Get block infos from... (in particular NFile, NDataPos and NUndoPos)

  pg := db.ConnectPg()
  defer pg.Close()

  res, err := db.GetHeaderFromHeight(pg, nHeight)
  if err != nil {
    log.Warn(err)
    return
  }
  btc.BlockHeader = res
  log.Info(fmt.Sprintf("%+v", btc))
  btc.getBlockFromFile()
  for _, tx := range(btc.Transactions) {
    _, err = db.InsertTransaction(pg, tx, btc.HashBlock)
    if err != nil {
      log.Warn(err)
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
  // btc.getTransaction() // Using index
}
