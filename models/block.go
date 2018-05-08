package models

import (
  "time"
)

type Script []byte
type Hash256 []byte

type TxInput struct {
  Hash              Hash256     `db:"hash"`
  Index             uint32      `db:"index"`
  Script            Script      `db:"script"`
  Sequence          uint32      `db:"sequence"`
  ScriptWitness     [][]byte
}

type TxOutput struct {
  Value             int64       `db:"value"`
  Script            Script      `db:"script"`
}

type Transaction struct {
  NVersion          int32       `db:"n_version"`
  Hash              Hash256     `db:"tx_hash"`
  Vin               []TxInput
  Vout              []TxOutput
  Locktime          uint32      `db:"locktime"`
	// StartPos      uint64 // not actually in blockchain data
}

// BlockIndexRecord contains general index records parameters
// It defines the structure of the postgres table
type BlockHeader struct {
  NVersion          int32     `db:"n_version"`    // Version
  NHeight           uint32    `db:"n_height"`             // (Index)
  NStatus           uint32    `db:"n_status"`// (Index)
  NTx               uint32    `db:"n_tx"`// (Index)
  NFile             uint32    `db:"n_file"`// (Index)
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
  BlockHeader
  Length        uint32
  Transactions  []Transaction
}
