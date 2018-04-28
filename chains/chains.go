package chains

import (
  "app/db"
  "github.com/jinzhu/gorm"
  "time"
  "fmt"
  "encoding/hex"
  // "github.com/syndtr/goleveldb/leveldb/opt"
)

type chains interface {
  getBlockHeader()
  // getBlockHeader()
  // getBlock()
}

type Script []byte

type TxInput struct {
	Hash          Hash256
	Index         uint32
	Script        Script
	Sequence      uint32
	ScriptWitness [][]byte
}

type TxOutput struct {
	Value  int64
	Script Script
}

type Transaction struct {
	hash     Hash256 // not actually in blockchain data; for caching
	Version  int32
	Locktime uint32
	Vin      []TxInput
	Vout     []TxOutput
	StartPos uint64 // not actually in blockchain data
}

// BlockIndexRecord contains general index records parameters
// It defines the structure of the postgres table
type BlockHeader struct {
  NVersion       int32    // Version
  NHeight        int32    `gorm:"primary_key;unique"`
  NStatus        uint32
  NTx            uint32
  NFile          int32
  NDataPos       uint32
  NUndoPos       uint32
  HashPrevBlock  Hash256  // previous block hash
  HashMerkleRoot Hash256  //
  NTime          time.Time  //
  NBits          uint32     //
  NNonce         uint32     //
  TargetDifficulty  uint32 // bits
}

type Block struct {
  gorm.Model
  BlockHeader
  Length        uint32
  Transactions  []Transaction `gorm:"-"` // don't store
}

func getBlockHeaders(c chains, nBlocks int) (bool, error) {

  // pg := db.Init()
  // pg.AutoMigrate(&BlockIndexRecord{})
  // defer pg.Close()
  for i := 0; i < nBlocks; i++ {
    c.getBlockHeader() // TODO: Errors checks
    // Copy, then insert in DB
    // pg.Create(&)
  }
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
  btcBlockHeader := &BtcBlockHeader{}
  // Reverse hex to get the LittleEndian order
  btcBlockHeader.HashPrevBlock = reverseHex(blockHashInBytes[:n])
  btcBlockHeader.IndexDb, _ = db.OpenIndexDb(btcDataDir) // TODO: Error handling
  defer btcBlockHeader.IndexDb.Close()
  // failIfReindexing(indexDb)
  getBlockHeaders(btcBlockHeader, 5)
}
