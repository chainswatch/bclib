package chains

import (
  "app/db"
  "fmt"
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
  pg := db.ConnectPg()
  err := pg.CreateTable(&Block{}).Error
  if err != nil {
    fmt.Println("Create Table Error:", err)
  }
  err = pg.CreateTable(&Transaction{}, &TxInput{}).Error
  if err != nil {
    fmt.Println("Create Table Error:", err)
  }
  fmt.Println("TABLES CREATED")
}

func ChainsWatcher() {
  //var qrlDataDir string = "./data/qrl/.qrl/data"

  btcBlock := &BtcBlock{}
  btcBlock.DataDir = "./chains/data/btc"

  startPg()

  // failIfReindexing(indexDb)
  // TODO: Init all tables. (db.AutoMigrate(&Model1{}, &Model2{}, &Model3{})
  // TODO: Use http://chainquery.com/bitcoin-api/getbestblockhash
  btcBlock.getBlockHashInBytes([]byte("000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"))
  BlockHeaderScanner(btcBlock, 5)
  BlockCoreScanner(btcBlock, 5)
}
