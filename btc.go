package main

import (
	"app/db"
  "log"
  "errors"
  "fmt"
	"app/generated/btc"
  "time"
  "encoding/hex"
)

// TODO: Maybe can optimize
func reverseHex(b []byte) []byte {
	newb := make([]byte, len(b))
	copy(newb, b)
	for i := len(newb)/2 - 1; i >= 0; i-- {
		opp := len(newb) - 1 - i
		newb[i], newb[opp] = newb[opp], newb[i]
	}

	return newb
}

func failIfReindexing(indexDb *db.IndexDb) {
	result, err := db.GetReindexing(indexDb)
	if err != nil {
		log.Fatal(err)
	}
	if result {
		log.Fatal(errors.New("bitcoind is reindexing"))
	}
}

type MagicId uint32

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func parseBlockHeaderFromFile(blockFile *BlockFile, block *btc.Block) error {
  var length = blockFile.ReadUint32() // TODO: Store it?
	block.Version = blockFile.ReadInt32()
	block.PrevBlock = blockFile.ReadBytes(32)
	block.MerkleRoot = blockFile.ReadBytes(32)
	block.Timestamp = int32(time.Unix(int64(blockFile.ReadUint32()), 0))
	block.Bits = blockFile.ReadUint32() // TODO: Parse this as mantissa?
	block.Nonce = blockFile.ReadUint32()

	return nil
}

func parseBlockTransactionFromFile(blockFile *BlockFile) (*btc.Transaction, error) {
	curPos, err := blockFile.Seek(0, 1)
	if err != nil {
		return nil, err
	}

	allowWitness := true // TODO: Port code - !(s.GetVersion() & SERIALIZE_TRANSACTION_NO_WITNESS);

	tx := &btc.Transaction{}
	// tx.StartPos = uint64(curPos)
	tx.Version = blockFile.ReadInt32()

	// Check for extended transaction serialization format
	p, _ := blockFile.Peek(1)
	var txInputLength uint64
	var txFlag byte
	if p[0] == 0 {
		// We are dealing with extended transaction
		blockFile.ReadByte()          // dummy
		txFlag = blockFile.ReadByte() // flags
		txInputLength = blockFile.ReadVarint()
	} else {
		txInputLength = blockFile.ReadVarint()
	}

	for i := uint64(0); i < txInputLength; i++ {
		input := TxInput{}
		input.Hash = blockFile.ReadBytes(32)
		input.Index = blockFile.ReadUint32() // TODO: Not sure if correctly read
		scriptLength := blockFile.ReadVarint()
		input.Script = blockFile.ReadBytes(scriptLength)
		input.Sequence = blockFile.ReadUint32()
		tx.Vin = append(tx.Vin, input)
	}

	txOutputLength := blockFile.ReadVarint()
	for i := uint64(0); i < txOutputLength; i++ {
		output := TxOutput{}
		output.Value = int64(blockFile.ReadUint64())
		scriptLength := blockFile.ReadVarint()
		output.Script = blockFile.ReadBytes(scriptLength)
		tx.Vout = append(tx.Vout, output)
	}

	if (txFlag&1) == 1 && allowWitness {
		txFlag ^= 1 // Not sure what this is for
		for i := uint64(0); i < txInputLength; i++ {
			witnessCount := blockFile.ReadVarint()
			tx.Vin[i].ScriptWitness = make([][]byte, witnessCount)
			for j := uint64(0); j < witnessCount; j++ {
				length := blockFile.ReadVarint()
				tx.Vin[i].ScriptWitness[j] = blockFile.ReadBytes(length)
			}
		}
	}

	tx.Locktime = blockFile.ReadUint32()

	return tx, nil
}

func parseBlockTransactionsFromFile(blockFile *BlockFile, block *btc.Block) error {
	// Read transaction count to know how many transactions to parse
  transactionCount := blockFile.ReadVarint() // TODO: Store it?
	//fmt.Printf("Total txns: %d\n", block.TransactionCount)
	for t := uint64(0); t < transactionCount; t++ {
		tx, err := ParseBlockTransactionFromFile(blockFile)
		if err != nil {
			return err
		}
		block.Transactions = append(block.Transactions, *tx)
	}

	return nil
}

func parseBlockFromFile(blockFile *BlockFile, magicHeader MagicId) (*btc.Block, error) {
	block := &btc.Block{}

	curPos, err := blockFile.Seek(0, 1)
	if err != nil {
		return nil, err
	}

	// Read and validate Magic ID
  magicId := MagicId(blockFile.ReadUint32())
	if magicId != magicHeader {
		blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
		return nil, errors.New("Invalid block header: Can't find Magic ID")
	}

	// Read header fields
	err = ParseBlockHeaderFromFile(blockFile, block)
	if err != nil {
		blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
		return nil, err
	}

	// Parse transactions
	err = ParseBlockTransactionsFromFile(blockFile, block)
	if err != nil {
		blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
		return nil, err
	}

	return block, nil
}

func newBlockFromFile(blockchainDataDir string, magicHeader MagicId, num uint32, pos uint32) (*btc.Block, error) {
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

  block, err := NewBlockFromFile(dataDir, uint32(result.NFile), result.NDataPos)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%+v\n", block)
  fmt.Printf("First Txid: %s\n", hex.EncodeToString(ReverseHex(block.Transactions[0].Txid())))
}

func btcWatcher(dataDir string) {
  indexDb, _ := db.OpenIndexDb(dataDir)
	defer indexDb.Close()

  failIfReindexing(indexDb)
  getBlock(indexDb, dataDir)
}
