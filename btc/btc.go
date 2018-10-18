package btc

/*

btc holds structs and methods used to parse the bitcoin blockchain

*/

import (
  "app/models"
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/jmoiron/sqlx"
  "os"
)

type MagicID uint32

// BtcBlockIndexRecord contains index records parameters specitic to BTC
type Btc struct {
  IndexDb       *leveldb.DB
  SqlDb         *sqlx.DB
	file          *os.File
  DataDir       string
  models.Block
}

type Tx struct {
	models.Transaction
}

func NewBtc(dataDir string) *Btc {
  return &Btc{DataDir: dataDir}
}

const (
  BLOCK_HAVE_DATA = 8  //!< full block available in blk*.dat
  BLOCK_HAVE_UNDO = 16 //!< undo data available in rev*.dat

  // Transaction types
	TX_UNKNOWN = 0x00
  TX_P2PKH = 0x01
  TX_P2SH = 0x02
  TX_P2PK = 0x03
  TX_MULTISIG = 0x04
  TX_P2WPKH = 0x05
  TX_P2WSH = 0x06

  TX_OPRETURN = 0x10 // Should contain data and not public key

  // Constants
  OP_0 = 0x00
  OP_1 = 0x51 // 1 is pushed
  OP_16 = 0x60

  OP_DUP = 0x76
  OP_HASH160 = 0xA9
  OP_CHECKSIG = 0xAC
  OP_PUSHDATA1 = 0x4C // Next byte containes the number of bytes to be pushed onto the stack
  OP_PUSHDATA2 = 0x4D // Next 2 bytes contain the number of bytes to be pushed (little endian)
  OP_PUSHDATA4 = 0x4E // Next 4 bytes contain the number of bytes to be pushed (little endian)

  // Bitwise logic
  OP_EQUAL = 0x87 // Returns 1 if the inputs are exactly equal, 0 otherwise
	OP_EQUALVERIFY = 0x88

  // Flow control
  OP_RETURN = 0x6A

  BTC_ECKEY_UNCOMPRESSED_LENGTH = 65
  BTC_ECKEY_COMPRESSED_LENGTH = 33
  // BTC_ECKEY_PKEY_LENGTH = 32
  // BTC_HASH_LENGTH = 32
	SHA256_DIGEST_LENGTH = 32
)
