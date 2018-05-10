package btc

import (
  "app/chains/parser"
  db "app/chains/repository"
  log "github.com/sirupsen/logrus"
  "github.com/syndtr/goleveldb/leveldb/util"
  "strings"
  "time"
)

// Parse raw data bytes
// https://github.com/bitcoin/bitcoin/blob/v0.15.1/src/chain.h#L387L407
func (btc *Btc) parseBlockHeaderData(data []byte) {
  // Parse the raw bytes
  dataBuf := parser.NewDataBuf(data)

  // Discard first varint
  // FIXME: Not exactly sure why need to, but if we don't do this we won't get correct values
  dataBuf.ShiftVarint()

  btc.NHeight = uint32(dataBuf.ShiftVarint())
  btc.NStatus = uint32(dataBuf.ShiftVarint())
  btc.NTx = uint32(dataBuf.ShiftVarint())
  if btc.NStatus & (BLOCK_HAVE_DATA|BLOCK_HAVE_UNDO) > 0 {
    btc.NFile = uint32(dataBuf.ShiftVarint())
  }
  if btc.NStatus & BLOCK_HAVE_DATA > 0 {
    btc.NDataPos = uint32(dataBuf.ShiftVarint())
  }
  if btc.NStatus & BLOCK_HAVE_UNDO > 0 {
    btc.NUndoPos = uint32(dataBuf.ShiftVarint())
  }

  btc.NVersion = dataBuf.Shift32bit()
  btc.HashPrevBlock = dataBuf.ShiftBytes(32)
  btc.HashMerkleRoot = dataBuf.ShiftBytes(32)
  btc.NTime = time.Unix(int64(dataBuf.ShiftU32bit()), 0)
  btc.NBits = dataBuf.ShiftU32bit()
  btc.NNonce = dataBuf.ShiftU32bit()
}

func (btc *Btc) GetBlockHeaders() {
  var err error
  btc.IndexDb, err = db.OpenIndexDb(btc.DataDir) // TODO: Error handling
  if err != nil {
    log.Warn("Error:", err)
  }
  defer btc.IndexDb.Close()

  iter := btc.IndexDb.NewIterator(util.BytesPrefix([]byte("b")), nil)
  for iter.Next() {
    btc.HashBlock = iter.Key()
    data := iter.Value()
    // log.Info(data)
    btc.parseBlockHeaderData(data)
    // Copy, then insert in DB
    _, err = db.InsertHeader(btc.SqlDb, btc.BlockHeader)
    if err != nil {
      if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
        log.Warn(err)
      }
    }
  }
  iter.Release()
  err = iter.Error()
  if err != nil {
    log.Warn(err)
  }
}

/*
func (btc *Btc) GetBlockHeader(nBlocks int) {
  var err error
  btc.IndexDb, err = db.OpenIndexDb(btc.DataDir) // TODO: Error handling
  if err != nil {
    log.Warn("Error:", err)
  }
  defer btc.IndexDb.Close()

  for i := 0; i < nBlocks; i++ {
    //fmt.Printf("Begin: blockHash: %v, %d bytes\n", btc.HashPrevBlock, len(btc.HashPrevBlock))

    // Get data
    data, err := btc.IndexDb.Get(append([]byte("b"), btc.HashPrevBlock...), nil)
    if err != nil {
      log.Warn(err)
    }
    btc.parseBlockHeaderData(data) // TODO: Errors checks
    // fmt.Printf("rawBlockHeader: %v\n", data)

    // Copy, then insert in DB
    _, err = db.InsertHeader(btc.SqlDb, btc.BlockHeader)
    if err != nil {
      if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
        log.Warn(err)
      }
    }
  }
  id, _ := db.GetRowCount(btc.SqlDb, "blocks")
  log.Info("blocks has ", id, " rows")
}
*/
