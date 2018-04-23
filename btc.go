package main

import (
	"app/db"
  "log"
  "errors"
  "fmt"
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

func NewBlockFromFile(blockchainDataDir string, magicHeader MagicId, num uint32, pos uint32) (*Block, error) {
	// Open file for reading
	blockFile, err := NewBlockFile(blockchainDataDir, num)
	if err != nil {
		return nil, err
	}
	defer blockFile.Close()

	// Seek to pos - 8 to start reading from block header
	fmt.Printf("Seeking to block at %d...\n", pos)
	_, err = blockFile.Seek(int64(pos-8), 0)
	if err != nil {
		return nil, err
	}

	return ParseBlockFromFile(blockFile, magicHeader)
}

func getBlock(indexDb *db.IndexDb, dataDir string) {
  result, err := db.GetBlockIndexRecordByBigEndianHex(indexDb, args[1])
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", result)

  block, err := NewBlockFromFile(dataDir, magicId, uint32(result.NFile), result.NDataPos)
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
  getBlock(&indexDb, dataDir)
}
