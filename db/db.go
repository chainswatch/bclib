package db

import (
  "fmt"
  "encoding/hex"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
  "time"
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

type IndexDb struct { // btc
	*leveldb.DB
}

type StateDb struct { // qrl
	*leveldb.DB
}

type Hash256 []byte
type MagicId uint32

type BlockIndexRecord struct {
	Version        int32
	Height         int32
	Status         uint32
	NTx            uint32
	NFile          int32
	NDataPos       uint32
	NUndoPos       uint32
	HashPrev       Hash256
	HashMerkleRoot Hash256
	NTime          time.Time
	NBits          uint32
	NNonce         uint32
}

func GetBlockIndexRecord(indexDb *IndexDb, blockHash []byte) (*BlockIndexRecord, error) {
	fmt.Printf("blockHash: %v, %d bytes\n", blockHash, len(blockHash))

	// Get data
	data, err := indexDb.Get(append([]byte("b"), blockHash...), nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("rawBlockIndexRecord: %v\n", data)

	// Parse the raw bytes
	blockIndexRecord := NewBlockIndexRecordFromBytes(data)

	return blockIndexRecord, nil
}

func NewBlockIndexRecordFromBytes(b []byte) *BlockIndexRecord {
	dataBuf := NewDataBuf(b)
	fmt.Printf("rawData: %v\n", b)
	dataHex := hex.EncodeToString(b)
	fmt.Printf("rawData: %v\n", dataHex)

	// Discard first varint
	// FIXME: Not exactly sure why need to, but if we don't do this we won't get correct values
	dataBuf.ShiftVarint()

	record := &BlockIndexRecord{}
	record.Height = int32(dataBuf.ShiftVarint())
	record.Status = uint32(dataBuf.ShiftVarint())
	record.NTx = uint32(dataBuf.ShiftVarint())
	if record.Status&(BLOCK_HAVE_DATA|BLOCK_HAVE_UNDO) > 0 {
		record.NFile = int32(dataBuf.ShiftVarint())
	}
	if record.Status&BLOCK_HAVE_DATA > 0 {
		record.NDataPos = uint32(dataBuf.ShiftVarint())
	}
	if record.Status&BLOCK_HAVE_UNDO > 0 {
		record.NUndoPos = uint32(dataBuf.ShiftVarint())
	}

	record.Version = dataBuf.Shift32bit()
	record.HashPrev = dataBuf.ShiftBytes(32)
	record.HashMerkleRoot = dataBuf.ShiftBytes(32)
	record.NTime = time.Unix(int64(dataBuf.ShiftU32bit()), 0)
	record.NBits = dataBuf.ShiftU32bit()
	record.NNonce = dataBuf.ShiftU32bit()

	return record
}

func OpenIndexDb(dataDir string) (*IndexDb, error) {
	db, err := leveldb.OpenFile(dataDir + "/blocks/index", &opt.Options{
		ReadOnly: true,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return &IndexDb{db}, nil
}

func OpenStateDb(dataDir string) (*StateDb, error) {
	db, err := leveldb.OpenFile(dataDir + "/chainstate", &opt.Options{
		ReadOnly: true,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return &StateDb{db}, nil
}

func GetReindexing(indexDb *IndexDb) (bool, error) {
	return indexDb.Has([]byte("R"), nil)
}

// TODO: Maybe can optimize
func reverseHex(b []byte) []byte {
	newb := make([]byte, len(b))
	copy(newb, b)
	for i := len(newb)/2 - 1; i >= 0; i-- {
		opp := len(newb) - 1 - i
		newb[i], newb[opp] = newb[opp], newb[i]
	}

	return newb
}

// Get block by block hash. TODO: By block height
func GetBlockIndexRecordByBigEndianHex(indexDb *IndexDb, blockHash string) (*BlockIndexRecord, error) {
	blockHashInBytes, err := hex.DecodeString(blockHash)
	if err != nil {
		return nil, err
	}
	// Reverse hex to get the LittleEndian order
	blockHashInBytes = reverseHex(blockHashInBytes)

	return GetBlockIndexRecord(indexDb, blockHashInBytes)
}
