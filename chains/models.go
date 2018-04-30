package chains

import (
  "time"
  "github.com/jinzhu/gorm"
)

// BlockIndexRecord contains general index records parameters
// It defines the structure of the postgres table
type BlockHeader struct {
  NVersion          int32    // Version
  NHeight           int32    `gorm:"primary_key;unique"`
  NStatus           uint32
  NTx               uint32
  NFile             int32
  NDataPos          uint32
  NUndoPos          uint32
  HashBlock         Hash256  `gorm:"-"` // current block hash
  HashPrevBlock     Hash256  // previous block hash
  HashMerkleRoot    Hash256  //
  NTime             time.Time  //
  NBits             uint32     //
  NNonce            uint32     //
  TargetDifficulty  uint32 // bits
}

type Block struct {
  gorm.Model
  BlockHeader
  Length        uint32
  Transactions  []Transaction `gorm:"-"` // don't store
}
