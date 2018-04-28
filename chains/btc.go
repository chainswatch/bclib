package chains

import (
  "github.com/syndtr/goleveldb/leveldb"
  "time"
  "fmt"
)

// btc holds structs and methods used to parse
// the bitcoin blockchain

type Hash256 []byte

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

/*
func GetBestBlock(chainstateDb *ChainstateDb) ([]byte, error) {
	return chainstateDb.Get([]byte("B"), nil)
}
*/

// BtcBlockIndexRecord contains index records parameters specitic to BTC
type BtcBlockHeader struct {
  IndexDb           *leveldb.DB
  hash              []byte
  BlockHeader
}

// Parse raw data bytes
// https://github.com/bitcoin/bitcoin/blob/v0.15.1/src/chain.h#L387L407
func (btc *BtcBlockHeader) getBlockHeader() {
  //fmt.Printf("Begin: blockHash: %v, %d bytes\n", btc.HashPrevBlock, len(btc.HashPrevBlock))

  // Get data
  data, err := btc.IndexDb.Get(append([]byte("b"), btc.HashPrevBlock...), nil)
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
  btc.HashBlock = append([]byte(nil), btc.HashPrevBlock...)
  btc.HashPrevBlock = dataBuf.ShiftBytes(32)
  btc.HashMerkleRoot = dataBuf.ShiftBytes(32)
  btc.NTime = time.Unix(int64(dataBuf.ShiftU32bit()), 0)
  btc.NBits = dataBuf.ShiftU32bit()
  btc.NNonce = dataBuf.ShiftU32bit()
  fmt.Printf("%+v\n", btc)
  //fmt.Printf("End: blockHash: %v, %d bytes\n", btc.HashPrevBlock, len(btc.HashPrevBlock))
}
