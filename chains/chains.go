package chains

import (
  db "app/chains/repository"
  "app/chains/btc"
)

type chains interface {
  GetBlockHeaders(int)
  GetBlock(int)
}

func BlockHeaderScanner(c chains, nBlocks int) {
  c.GetBlockHeaders(nBlocks)
}

func BlockCoreScanner(c chains, nHeight int) {
  c.GetBlock(nHeight)
}

// startPg starts the postgres migration process
func startPg() {
	db.CreateBtc()
}

func ChainsWatcher() {
  //var qrlDataDir string = "./data/qrl/.qrl/data"

  btc := btc.NewBtc("/data/chains/btc")

  startPg()

  // failIfReindexing(indexDb)
  // TODO: Init all tables. (db.AutoMigrate(&Model1{}, &Model2{}, &Model3{})
  // TODO: Use http://chainquery.com/bitcoin-api/getbestblockhash
  btc.GetBlockHashInBytes([]byte("0000000026f34d197f653c5e80cb805e40612eadb0f45d00d7ea4164a20faa33"))
  BlockHeaderScanner(btc, 51)
  BlockCoreScanner(btc, 1)
}
