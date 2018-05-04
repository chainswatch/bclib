package chains

import (
  "app/models"
  db "app/chains/repository"
  "fmt"
)

// Witness : https://github.com/bitcoin/bitcoin/blob/master/src/primitives/transaction.h
const SERIALIZE_TRANSACTION_NO_WITNESS = 0x40000000;

/*
* Basic transaction serialization format:
* - int32_t nVersion
* - std::vector<CTxIn> vin
* - std::vector<CTxOut> vout
* - uint32_t nLockTime
*
* Extended transaction serialization format:
* - int32_t nVersion
* - unsigned char dummy = 0x00
* - unsigned char flags (!= 0)
* - std::vector<CTxIn> vin
* - std::vector<CTxOut> vout
* - if (flags & 1):
*   - CTxWitness wit;
* - uint32_t nLockTime
*/

func (btc *BtcBlock) parseBlockTransactionFromFile() (*models.Transaction, error) {
  // curPos, err := btc.Seek(0, 1)
  allowWitness := true // TODO: Port code - !(s.GetVersion() & SERIALIZE_TRANSACTION_NO_WITNESS);

  tx := &models.Transaction{}
  // tx.StartPos = uint64(curPos)
  tx.NVersion = btc.ReadInt32() // 

  // Check for extended transaction serialization format
  var txInputLength uint64
  var txFlag byte
  // Try to read. Look for dummy
  p, _ := btc.Peek(1)
  if p[0] == 0 {
    // We are dealing with extended transaction
    btc.ReadByte()          // dummy (0x00)
    txFlag = btc.ReadByte() // flags (!=0)
    if (txFlag != 0) {
      txInputLength = btc.ReadVarint()
    }
  } else {
    txInputLength = btc.ReadVarint()
  }

  for i := uint64(0); i < txInputLength; i++ {
    input := models.TxInput{}
    input.Hash = btc.ReadBytes(32)
    input.Index = btc.ReadUint32() // TODO: Not sure if correctly read
    scriptLength := btc.ReadVarint()
    input.Script = btc.ReadBytes(scriptLength)
    input.Sequence = btc.ReadUint32()
    tx.TxInputs = append(tx.TxInputs, input)
  }

  txOutputLength := btc.ReadVarint()
  for i := uint64(0); i < txOutputLength; i++ {
    output := models.TxOutput{}
    output.Value = int64(btc.ReadUint64())
    scriptLength := btc.ReadVarint()
    output.Script = btc.ReadBytes(scriptLength)
    tx.TxOutputs = append(tx.TxOutputs, output)
  }

  if (txFlag & 1) == 1 && allowWitness {
    txFlag ^= 1 // Not sure what this is for
    for i := uint64(0); i < txInputLength; i++ {
      witnessCount := btc.ReadVarint()
      tx.TxInputs[i].ScriptWitness = make([][]byte, witnessCount)
      for j := uint64(0); j < witnessCount; j++ {
        length := btc.ReadVarint()
        tx.TxInputs[i].ScriptWitness[j] = btc.ReadBytes(length)
      }
    }
  }

  tx.Locktime = btc.ReadUint32()

  return tx, nil
}

func (btc *BtcBlock) getTransaction() {
  f, _ := db.GetFlag(btc.IndexDb, []byte("txindex"))
  if !f {
    fmt.Println("txindex is not enabled for your bitcoind")
  }
  /*
  result, err := db.GetTxIndexRecordByBigEndianHex(indexDb, args[1])
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", result)

  tx, err := blockchainparser.NewTxFromFile(datadir, magicId, uint32(result.NFile), result.NDataPos, result.NTxOffset)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%+v\n", tx)
  */
}

func (btc *BtcBlock) parseBlockTransactionsFromFile() error {
  // Read transaction count to know how many transactions to parse
  TransactionCount := btc.ReadVarint()
  // TODO: TransactionCount : Compare with Index stored value (NTx)
  fmt.Printf("Total txns: %d vs %d\n", TransactionCount, btc.NTx)
  for t := uint64(0); t < TransactionCount; t++ {
    tx, err := btc.parseBlockTransactionFromFile()
    if err != nil {
      return err
    }
    fmt.Printf("Transaction: %+v\n", tx)
    btc.Transactions = append(btc.Transactions, *tx)
  }
  btc.saveTransactions()
  return nil
}

func (btc *BtcBlock) saveTransactions() {
  /*
  btc.IndexDb, _ = db.OpenIndexDb(btc.DataDir) // TODO: Error handling
  defer btc.IndexDb.Close()

  pg := db.PgInit()
  pg.AutoMigrate(&TxInput{})
  pg.AutoMigrate(&Transaction{})
  err := pg.Create(&btc.Transactions)
  if err != nil {
    fmt.Println(err) // 
  }
  */
}


