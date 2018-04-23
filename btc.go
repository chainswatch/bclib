package main

import (
	"app/db"
  "log"
  "errors"
)

func failIfReindexing(indexDb *db.IndexDb) {
	result, err := db.GetReindexing(indexDb)
	if err != nil {
		log.Fatal(err)
	}
	if result {
		log.Fatal(errors.New("bitcoind is reindexing"))
	}
}

func getBlock() {
  result, err := db.GetBlockIndexRecordByBigEndianHex(indexDb, args[1])
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", result)

  block, err := blockchainparser.NewBlockFromFile(datadir, magicId, uint32(result.NFile), result.NDataPos)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%+v\n", block)
  fmt.Printf("First Txid: %s\n", hex.EncodeToString(blockchainparser.ReverseHex(block.Transactions[0].Txid())))
}

func btcWatcher(dataDir string) {
  indexDb, _ := db.OpenIndexDb(dataDir)
	defer indexDb.Close()

  failIfReindexing(indexDb)
  getBlock()
}
