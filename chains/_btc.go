// +build ignore
package chains

import (
  "app/db"
  //"app/chains"
  "log"
  "errors"
  "fmt"
  "encoding/hex"
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

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func parseBlockHeaderFromFile(blockFile *BlockFile, block *btc.BlockHeader) error {
  var length = blockFile.ReadUint32() // TODO: Store it?
  fmt.Println("BlockHeader length", length)
  block.nVersion = blockFile.ReadInt32()
  block.hashPrev = blockFile.ReadBytes(32)
  block.hashMerkleRoot = blockFile.ReadBytes(32)
  block.nTime = blockFile.ReadUint32()
  block.nBits = blockFile.ReadUint32() // TODO: Parse this as mantissa?
  block.nNonce = blockFile.ReadUint32()

  return nil
}

func parseBlockTransactionFromFile(blockFile *BlockFile) (*btc.Transaction, error) {
  curPos, err := blockFile.Seek(0, 1)
  fmt.Println("CurPos:", curPos)
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
    input := &btc.TxInput{}
    input.OutpointHash = blockFile.ReadBytes(32)
    input.OutpointIndex = blockFile.ReadUint32() // TODO: Not sure if correctly read
    scriptLength := blockFile.ReadVarint()
    fmt.Println("scriptLength:", scriptLength)
    // input.SigScript = blockFile.ReadBytes(scriptLength) // TODO
    // input.Sequence = blockFile.ReadUint32()
    tx.Inputs = append(tx.Inputs, input)
  }

  txOutputLength := blockFile.ReadVarint()
  for i := uint64(0); i < txOutputLength; i++ {
    output := &btc.TxOutput{}
    output.Value = uint32(blockFile.ReadUint64())
    scriptLength := blockFile.ReadVarint()
    fmt.Println("scriptLength:", scriptLength)
    // output.Script = blockFile.ReadBytes(scriptLength)
    tx.Outputs = append(tx.Outputs, output)
  }

  if (txFlag&1) == 1 && allowWitness {
    txFlag ^= 1 // Not sure what this is for
    for i := uint64(0); i < txInputLength; i++ {
      witnessCount := blockFile.ReadVarint()
      // tx.Inputs[i].ScriptWitness = make([][]byte, witnessCount)
      for j := uint64(0); j < witnessCount; j++ {
        length := blockFile.ReadVarint()
        fmt.Println("Length:", length)
        // tx.Inputs[i].ScriptWitness[j] = blockFile.ReadBytes(length)
      }
    }
  }

  tx.LockTime = blockFile.ReadUint32()

  return tx, nil
}

func parseBlockTransactionsFromFile(blockFile *BlockFile, block *btc.Block) error {
  // Read transaction count to know how many transactions to parse
  transactionCount := blockFile.ReadVarint() // TODO: Store it?
  //fmt.Printf("Total txns: %d\n", block.TransactionCount)
  for t := uint64(0); t < transactionCount; t++ {
    tx, err := parseBlockTransactionFromFile(blockFile)
    if err != nil {
      return err
    }
    block.Transactions = append(block.Transactions, tx)
  }
  return nil
}

func parseBlockFromFile(blockFile *BlockFile) (*btc.Block, error) {
  block := &btc.Block{}

  curPos, err := blockFile.Seek(0, 1)
  if err != nil {
    return nil, err
  }

  // Read and validate Magic ID
  magicID := uint32(blockFile.ReadUint32())
  fmt.Println("Magic ID:", magicID)
  /*
  if magicId != magicHeader {
    blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    return nil, errors.New("Invalid block header: Can't find Magic ID")
  }
  */

  // Read header fields
  err = parseBlockHeaderFromFile(blockFile, block)
  if err != nil {
    blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    return nil, err
  }

  // Parse transactions
  err = parseBlockTransactionsFromFile(blockFile, block)
  if err != nil {
    blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    return nil, err
  }

  return block, nil
}

func newBlockFromFile(blockchainDataDir string, num uint32, pos uint32) (*btc.Block, error) {
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

  return parseBlockFromFile(blockFile)
}

func getBlock(dataDir string, nFile uint32, nDataPos uint32) {
  block, err := newBlockFromFile(dataDir, nFile, nDataPos)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%+v\n", block)
  // fmt.Printf("First Txid: %s\n", hex.EncodeToString(reverseHex(block.Transactions[0].Txid())))
}

func getBlockHeaders(indexDb *db.IndexDb, blockHash []byte, nBlocks int) (bool, error) {
  var header *db.BlockHeader

  blockHashInBytes := make([]byte, hex.DecodedLen(len(blockHash)))
  n, err := hex.Decode(blockHashInBytes, blockHash)
  if err != nil {
    return true, err
  }
  // Reverse hex to get the LittleEndian order
  // blockHashInBytes = reverseHex(blockHashInBytes)
  blockHashInBytes = reverseHex(blockHashInBytes[:n])

  pg := db.DbInit()
  pg.AutoMigrate(&db.BlockHeader{})
  defer pg.Close()
  for i := 0; i < nBlocks; i++ {
    if i == 0 {
      header, err = db.GetBlockHeader(indexDb, blockHashInBytes)
    } else {
      header, err = db.GetBlockHeader(indexDb, header.hashPrev)
    }
    if err != nil {
      log.Fatal(err)
    }
    pg.Create(&header)
    fmt.Printf("%+v\n", header)
  }
  var count int
  pg.Table("block_headers").Count(&count)
  fmt.Println(count, "RECORDS")
  return true, nil
}

func btcWatcher(dataDir string) {
  indexDb, _ := db.OpenIndexDb(dataDir)
  defer indexDb.Close()

  failIfReindexing(indexDb)
  getBlockHeaders(indexDb, []byte("000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"), 1)
}
