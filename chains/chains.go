package chains

import (
  db "app/chains/repository"
)

type chains interface {
  getBlockHeaders(int)
  getBlock(int)
}

func BlockHeaderScanner(c chains, nBlocks int) {
  c.getBlockHeaders(nBlocks)
}

func BlockCoreScanner(c chains, nHeight int) {
  c.getBlock(nHeight)
}

// startPg starts the postgres migration process
func startPg() {
	db.CreateBtc()
}

func ChainsWatcher() {
  //var qrlDataDir string = "./data/qrl/.qrl/data"

  btcBlock := &BtcBlock{}
  btcBlock.DataDir = "/data/chains/btc"

  startPg()

  // failIfReindexing(indexDb)
  // TODO: Init all tables. (db.AutoMigrate(&Model1{}, &Model2{}, &Model3{})
  // TODO: Use http://chainquery.com/bitcoin-api/getbestblockhash
  btcBlock.getBlockHashInBytes([]byte("000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"))
  BlockHeaderScanner(btcBlock, 5)
  btcBlock2 := &BtcBlock{}
  btcBlock2.DataDir = "/data/chains/btc"
  BlockCoreScanner(btcBlock2, 5)
}
