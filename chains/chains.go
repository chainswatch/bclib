package chains

type chains interface {
  getBlockHeaders(int)
  // getBlockHeader()
  // getBlock()
}

type Script []byte

type TxInput struct {
	Hash          Hash256
	Index         uint32
	Script        Script
	Sequence      uint32
	ScriptWitness [][]byte
}

type TxOutput struct {
	Value  int64
	Script Script
}

type Transaction struct {
	Hash          Hash256 // not actually in blockchain data; for caching
	Version       int32
	Locktime      uint32
	Vin           []TxInput
	Vout          []TxOutput
	StartPos      uint64 // not actually in blockchain data
}

// TODO: Rename
func BlockHeaderScanner(c chains, nBlocks int) {
  c.getBlockHeaders(nBlocks)
}

func ChainsWatcher() {
  //var qrlDataDir string = "./data/qrl/.qrl/data"

  btcBlock := &BtcBlock{}
  btcBlock.DataDir = "./chains/data/btc"

  // failIfReindexing(indexDb)
  btcBlock.getBlockHashInBytes([]byte("000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"))
  BlockHeaderScanner(btcBlock, 5)
}
