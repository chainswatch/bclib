package chains

import (
  "app/misc"
  db "app/chains/repository"
  log "github.com/sirupsen/logrus"
  "app/chains/btc"
)

type chains interface {
  // GetBlockHeader(int)
  GetBlockHeaders()
  GetAllBlocks()
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

func ChainsWatcher() {
  //var qrlDataDir string = "./data/qrl/.qrl/data"

  btc := btc.NewBtc("/data/chains/btc")

	btc.SqlDb = db.Create(true)

  // failIfReindexing(indexDb)
  // TODO: Init all tables. (db.AutoMigrate(&Model1{}, &Model2{}, &Model3{})
  // TODO: Use http://chainquery.com/bitcoin-api/getbestblockhash
  // log.SetLevel(log.WarnLevel)

  log.Info("START")
  btc.HashPrevBlock = misc.HexToByte([]byte("00000000dfd5d65c9d8561b4b8f60a63018fe3933ecb131fb37f905f87da951a"))
  BlockCoreScanner(btc)
  log.Info("DONE")
  //BlockHeaderScanner(btc)
  //BlockCoreScanner(btc, 2000)
  //log.Info("Core DONE")
  /*
  btc.HashPrevBlock = misc.HexToByte([]byte("0000000026f34d197f653c5e80cb805e40612eadb0f45d00d7ea4164a20faa33"))
  BlockHeaderScanner(btc, 51)
  log.Info("Header DONE")
  log.Info("Core DONE")
  */
}
