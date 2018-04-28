package chains

import (
  "app/db"
  "github.com/jinzhu/gorm"
  "time"
  "github.com/syndtr/goleveldb/leveldb"
  "fmt"
  "encoding/hex"
  // "github.com/syndtr/goleveldb/leveldb/opt"
)

const (
  //! Unused.
  BLOCK_VALID_UNKNOWN = 0

  //! Parsed, version ok, hash satisfies claimed PoW, 1 <= vtx count <= max, timestamp not in future
  BLOCK_VALID_HEADER = 1

  //! All parent headers found, difficulty matches, timestamp >= median previous, checkpoint. Implies all parents
  //! are also at least TREE.
  BLOCK_VALID_TREE = 2

  /**
  * Only first tx is coinbase, 2 <= coinbase input script length <= 100, transactions valid, no duplicate txids,
  * sigops, size, merkle root. Implies all parents are at least TREE but not necessarily TRANSACTIONS. When all
  * parent blocks also have TRANSACTIONS, CBlockIndex::nChainTx will be set.
  */
  BLOCK_VALID_TRANSACTIONS = 3

  //! Outputs do not overspend inputs, no double spends, coinbase output ok, no immature coinbase spends, BIP30.
  //! Implies all parents are also at least CHAIN.
  BLOCK_VALID_CHAIN = 4

  //! Scripts & signatures ok. Implies all parents are also at least SCRIPTS.
  BLOCK_VALID_SCRIPTS = 5

  //! All validity bits.
  BLOCK_VALID_MASK = BLOCK_VALID_HEADER | BLOCK_VALID_TREE | BLOCK_VALID_TRANSACTIONS |
  BLOCK_VALID_CHAIN | BLOCK_VALID_SCRIPTS

  BLOCK_HAVE_DATA = 8  //!< full block available in blk*.dat
  BLOCK_HAVE_UNDO = 16 //!< undo data available in rev*.dat
  BLOCK_HAVE_MASK = BLOCK_HAVE_DATA | BLOCK_HAVE_UNDO

  BLOCK_FAILED_VALID = 32 //!< stage after last reached validness failed
  BLOCK_FAILED_CHILD = 64 //!< descends from failed block
  BLOCK_FAILED_MASK  = BLOCK_FAILED_VALID | BLOCK_FAILED_CHILD

  BLOCK_OPT_WITNESS = 128 //!< block data in blk*.data was received with a witness-enforcing client
)


type chains interface {
  getBlockIndexRecord()
  // getBlockHeader()
  // getBlock()
}

// BtcBlockIndexRecord contains index records parameters specitic to BTC
type BtcBlockIndexRecord struct {
  IndexDb        *leveldb.DB
  NVersion       int32
  NHeight        int32
  NStatus        uint32
  NTx            uint32
  NFile          int32
  NDataPos       uint32
  NUndoPos       uint32
  HashPrev       []byte
  HashMerkleRoot []byte
  NTime          time.Time
  NBits          uint32
  NNonce         uint32
}

// BlockIndexRecord contains general index records parameters
// It defines the structure of the postgres table
type BlockIndexRecord struct {
  gorm.Model
  NVersion       int32
  NHeight        int32    `gorm:"primary_key;unique"`
  NStatus        uint32
  NTx            uint32
  NFile          int32
  NDataPos       uint32
  NUndoPos       uint32
  HashPrev       []byte
  HashMerkleRoot []byte
  NTime          time.Time
  NBits          uint32
  NNonce         uint32
}

// Parse raw data bytes
// https://github.com/bitcoin/bitcoin/blob/v0.15.1/src/chain.h#L387L407
func (btc BtcBlockIndexRecord) getBlockIndexRecord() {
  // fmt.Printf("blockHash: %v, %d bytes\n", blockHash, len(blockHash))

  // Get data
  data, err := btc.IndexDb.Get(append([]byte("b"), btc.HashPrev...), nil)
  if err != nil {
    fmt.Printf("Error")
  }
  // fmt.Printf("rawBlockHeader: %v\n", data)

  // Parse the raw bytes
  dataBuf := NewDataBuf(data)
  //fmt.Printf("rawData: %v\n", b)
  //dataHex := hex.EncodeToString(b)
  //fmt.Printf("rawData: %v\n", dataHex)

  // Discard first varint
  // FIXME: Not exactly sure why need to, but if we don't do this we won't get correct values
  dataBuf.ShiftVarint()

  btc.NHeight = int32(dataBuf.ShiftVarint())
  btc.NStatus = uint32(dataBuf.ShiftVarint())
  btc.NTx = uint32(dataBuf.ShiftVarint())
  if btc.NStatus & (BLOCK_HAVE_DATA|BLOCK_HAVE_UNDO) > 0 {
    btc.NFile = int32(dataBuf.ShiftVarint())
  }
  if btc.NStatus & BLOCK_HAVE_DATA > 0 {
    btc.NDataPos = uint32(dataBuf.ShiftVarint())
  }
  if btc.NStatus & BLOCK_HAVE_UNDO > 0 {
    btc.NUndoPos = uint32(dataBuf.ShiftVarint())
  }

  btc.NVersion = dataBuf.Shift32bit()
  btc.HashPrev = dataBuf.ShiftBytes(32)
  btc.HashMerkleRoot = dataBuf.ShiftBytes(32)
  btc.NTime = time.Unix(int64(dataBuf.ShiftU32bit()), 0)
  btc.NBits = dataBuf.ShiftU32bit()
  btc.NNonce = dataBuf.ShiftU32bit()
  fmt.Printf("%+v\n", btc)
}

func getBlockIndexRecords(c chains, nBlocks int) (bool, error) {

  // pg := db.Init()
  // pg.AutoMigrate(&BlockIndexRecord{})
  // defer pg.Close()
  for i := 0; i < nBlocks; i++ {
    c.getBlockIndexRecord() // TODO: Errors checks
    // Copy, then insert in DB
    // pg.Create(&)
    // fmt.Printf("%+v\n", c...)
  }
  // var count int
  // pg.Table("block_headers").Count(&count)
  // fmt.Println(count, "RECORDS")
  return true, nil
}

func ChainsWatcher() {
  //var qrlDataDir string = "./data/qrl/.qrl/data"

  var btcDataDir = "./chains/data/btc"

  blockHashStart := []byte("000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506")
  blockHashInBytes := make([]byte, hex.DecodedLen(len(blockHashStart)))
  n, err := hex.Decode(blockHashInBytes, blockHashStart)
  if err != nil {
    fmt.Println(err)
  }
  btcIndexRecord := &BtcBlockIndexRecord{}
  // Reverse hex to get the LittleEndian order
  btcIndexRecord.HashPrev = reverseHex(blockHashInBytes[:n])
  btcIndexRecord.IndexDb, _ = db.OpenIndexDb(btcDataDir) // TODO: Error handling
  defer btcIndexRecord.IndexDb.Close()
  // failIfReindexing(indexDb)
  getBlockIndexRecords(btcIndexRecord, 5)
}
