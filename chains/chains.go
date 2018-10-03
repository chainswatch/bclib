package chains

import (
  //"app/misc"
  db "app/chains/repository"
  log "github.com/sirupsen/logrus"
  "app/chains/btc"
)

type chains interface {
  // GetBlockHeader(int)
  GetBlockHeaders()
  GetAllBlocks()
  MempoolWatcher()
}

func BlockHeaderScanner(c chains) {
  c.GetBlockHeaders()
}

/*
func BlockHeaderGetter(c chains, nBlocks int) {
  c.GetBlockHeader(nBlocks)
}
*/

func BlockCoreScanner(c chains) {
  c.GetAllBlocks()
}

func MempoolScanner(c chains) {
  c.MempoolWatcher()
}

func ChainsWatcher() {
  log.SetLevel(log.DebugLevel)
  //log.SetLevel(log.InfoLevel)
  //var qrlDataDir string = "./data/qrl/.qrl/data"

  btc := btc.NewBtc("/data/crypto/chains/btc")

  btc.SqlDb = db.Create(false)

  // failIfReindexing(indexDb)
  // TODO: Init all tables. (db.AutoMigrate(&Model1{}, &Model2{}, &Model3{})
  // TODO: Use http://chainquery.com/bitcoin-api/getbestblockhash

  // Load blockchain history
  log.Info("BTC Transactions listener: Start")
  MempoolScanner(btc)
  log.Info("BTC Transactions listener: Stop")

  // TODO: Test if connected to database
  log.Info("BTC Blockchain Sync: Start")
  //btc.HashPrevBlock = misc.HexToByte([]byte("00000000dfd5d65c9d8561b4b8f60a63018fe3933ecb131fb37f905f87da951a"))
  //BlockCoreScanner(btc)
  log.Info("BTC Blockchain Sync: Finished")
}
