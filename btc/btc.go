package btc

/*
btc holds structs and methods used to parse the bitcoin blockchain
*/

type magicID uint32

// BtcBlockIndexRecord contains index records parameters specitic to BTC
const (
	blockHaveData = 8  //!< full block available in blk*.dat
	blockHaveUndo = 16 //!< undo data available in rev*.dat

	txP2pkh    = 0x01
	txP2sh     = 0x02
	txP2pk     = 0x03
	txMultisig = 0x04
	txP2wpkh   = 0x05
	txP2wsh    = 0x06 // bench32

	txOpreturn = 0x10 // Should contain data and not public key
	txParseErr = 0xfe
	txUnknown  = 0xff

	op0  = 0x00
	op1  = 0x51 // 1 is pushed
	op16 = 0x60

	opDup       = 0x76
	opHash160   = 0xA9
	opChecksig  = 0xAC
	opPushdata1 = 0x4C // Next byte containes the number of bytes to be pushed onto the stack
	opPushdata2 = 0x4D // Next 2 bytes contain the number of bytes to be pushed (little endian)
	opPushdata4 = 0x4E // Next 4 bytes contain the number of bytes to be pushed (little endian)

	opEqual       = 0x87 // Returns 1 if the inputs are exactly equal, 0 otherwise
	opEqualverify = 0x88

	opReturn = 0x6A

	btcEckeyUncompressedLength = 65
	btcEckeyCompressedLength   = 33
	sha256DigestLength         = 32

	// BTC_ECKEY_PKEY_LENGTH = 32
	// BTC_HASH_LENGTH = 32
)
