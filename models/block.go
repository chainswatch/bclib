package models

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
  NVersion          int32     `db:"n_version"`    // Version
  NHeight           int32     `db:"n_height"`             // (Index)
  NStatus           uint32    `db:"n_status"`// (Index)
  NTx               uint32    `db:"n_tx"`// (Index)
  NFile             int32     `db:"n_file"`// (Index)
  NDataPos          uint32    `db:"n_data_pos"`// (Index)
  NUndoPos          uint32    `db:"n_undo_pos"`// (Index)
  HashBlock         Hash256   `db:"hash_block"`// current block hash (Added)
  HashPrevBlock     Hash256   `db:"hash_prev_block"`// previous block hash (Index)
  HashMerkleRoot    Hash256   `db:"hash_merkle_root"`// (Index)
  NTime             time.Time `db:"n_time"`// (Index)
  NBits             uint32    `db:"n_bits"`// (Index)
  NNonce            uint32    `db:"n_nonce"`// (Index)
  TargetDifficulty  uint32    `db:"target_difficulty"`//
}

type Block struct {
  CreatedAt     time.Time
  UpdatedAt     time.Time
  DeletedAt     *time.Time
  BlockHeader
  Length        uint32
  Transactions  []Transaction `gorm:"-"` // don't store
}
