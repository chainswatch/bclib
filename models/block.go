package models

// TxInput holds tx inputs
type TxInput struct {
	Hash          []byte `db:"hash"`     // Hash previous tx
	Index         uint32 `db:"index"`    // Output previous tx
	Script        []byte `db:"script"`   // Useless?
	Sequence      uint32 `db:"sequence"` // Always 0xFFFFFFFF
	ScriptWitness [][]byte
}

// TxOutput holds tx outputs
type TxOutput struct {
	Index uint32 `db:"index"` // Output index
	Value uint64 `db:"value"` // Satoshis
	// TODO: Add type
	Addr     []byte `db:"addr"` // Public key
	AddrType uint8
	Script   []byte `db:"script"` // Where the magic happens
}

// Tx holds transaction
type Tx struct {
	NVersion int32  `db:"n_version"` // Always 1 or 2
	Hash     []byte `db:"tx_hash"`   // Transaction hash (computed)
	NVin     uint32 `db:"n_vin"`     // Number of inputs
	NVout    uint32 `db:"n_vout"`    // Number of outputs
	Vin      []TxInput
	Vout     []TxOutput
	Locktime uint32 `db:"locktime"`
	Segwit		bool
}

// BlockHeader contains general index records parameters
// It defines the structure of the postgres table
type BlockHeader struct {
	NVersion         uint32 `db:"n_version"`         // Version
	NHeight          uint32 `db:"n_height"`          //
	NStatus          uint32 `db:"n_status"`          // ???
	NTx              uint32 `db:"n_tx"`              // Number of txs
	NFile            uint32 `db:"n_file"`            // File number
	NDataPos         uint32 `db:"n_data_pos"`        // (Index)
	NUndoPos         uint32 `db:"n_undo_pos"`        // (Index)
	Hash             []byte `db:"hash_block"`        // current block hash (Added)
	HashPrev         []byte `db:"hash_prev_block"`   // previous block hash (Index)
	HashMerkleRoot   []byte `db:"hash_merkle_root"`  //
	NTime            uint32 `db:"n_time"`            // (Index)
	NBits            uint32 `db:"n_bits"`            // (Index)
	NNonce           uint32 `db:"n_nonce"`           // (Index)
	TargetDifficulty uint32 `db:"target_difficulty"` //
	NSize            uint32 `db:"n_size"`            // Block size
}

// Block contains block infos
type Block struct {
	BlockHeader
	Txs []Tx
}
