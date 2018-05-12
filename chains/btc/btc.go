package btc

import (
  "app/models"
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/jmoiron/sqlx"
  "os"
)

// btc holds structs and methods used to parse
// the bitcoin blockchain

type MagicID uint32

// BtcBlockIndexRecord contains index records parameters specitic to BTC
type Btc struct {
  IndexDb       *leveldb.DB
  SqlDb         *sqlx.DB
	file          *os.File
  DataDir       string
  models.Block
}

func NewBtc(dataDir string) *Btc {
  return &Btc{DataDir: dataDir}
}

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
