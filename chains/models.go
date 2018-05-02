package chains 

import (
  "time"
  // "github.com/jinzhu/gorm"
)

type Script []byte
type Hash256 []byte

type TxInput struct {
  TransactionID     uint
  Hash              Hash256     `gorm:"-"`
	Index             uint32
	Script            Script
	Sequence          uint32
  ScriptWitness     [][]byte    `gorm:"-"`
}

type TxOutput struct {
	Value             int64
	Script            Script
}

type Transaction struct {
	NVersion          int32
  TxInputs          []TxInput
  TxOutputs         []TxOutput  `gorm:"-"`
	Locktime          uint32
	// StartPos      uint64 // not actually in blockchain data
}

// BlockIndexRecord contains general index records parameters
// It defines the structure of the postgres table
type BlockHeader struct {
  NVersion          int32    // Version
  NHeight           int32    `gorm:"primary_key"` // NOT NULL & UNIQUE (TODO: Combination primary key)
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
  CreatedAt     time.Time
  UpdatedAt     time.Time
  DeletedAt     *time.Time
  BlockHeader
  Length        uint32
  Transactions  []Transaction `gorm:"-"` // don't store
}
