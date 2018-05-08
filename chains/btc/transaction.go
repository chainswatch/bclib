package btc

import (
  "app/models"
  "app/chains/parser"
  db "app/chains/repository"
  log "github.com/sirupsen/logrus"
  "encoding/binary"
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

func (btc *Btc) parseBlockTransactionFromFile(blockFile *parser.BlockFile) (*models.Transaction, error) {
  // curPos, err := btc.Seek(0, 1)
  allowWitness := true // TODO: Port code - !(s.GetVersion() & SERIALIZE_TRANSACTION_NO_WITNESS);

  tx := &models.Transaction{}
  // tx.StartPos = uint64(curPos)
  tx.NVersion = blockFile.ReadInt32() // 

  // Check for extended transaction serialization format
  var txInputLength uint64
  var txFlag byte
  // Try to read. Look for dummy
  p, _ := blockFile.Peek(1)
  if p[0] == 0 {
    // We are dealing with extended transaction
    blockFile.ReadByte()          // dummy (0x00)
    txFlag = blockFile.ReadByte() // flags (!=0)
    if (txFlag != 0) {
      txInputLength = blockFile.ReadVarint()
    }
  } else {
    txInputLength = blockFile.ReadVarint()
  }

  for i := uint64(0); i < txInputLength; i++ {
    input := models.TxInput{}
    input.Hash = blockFile.ReadBytes(32)
    input.Index = blockFile.ReadUint32() // TODO: Not sure if correctly read
    scriptLength := blockFile.ReadVarint()
    input.Script = blockFile.ReadBytes(scriptLength)
    input.Sequence = blockFile.ReadUint32()
    tx.Vin = append(tx.Vin, input)
  }

  txOutputLength := blockFile.ReadVarint()
  for i := uint64(0); i < txOutputLength; i++ {
    output := models.TxOutput{}
    output.Value = int64(blockFile.ReadUint64())
    scriptLength := blockFile.ReadVarint()
    output.Script = blockFile.ReadBytes(scriptLength)
    tx.Vout = append(tx.Vout, output)
  }

  if (txFlag & 1) == 1 && allowWitness {
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

func varint(n uint64) []byte {
	if n > 4294967295 {
		val := make([]byte, 8)
		binary.BigEndian.PutUint64(val, n)
		return append([]byte{0xFF}, val...)
	} else if n > 65535 {
		val := make([]byte, 4)
		binary.BigEndian.PutUint32(val, uint32(n))
		return append([]byte{0xFE}, val...)
	} else if n > 255 {
		val := make([]byte, 2)
		binary.BigEndian.PutUint16(val, uint16(n))
		return append([]byte{0xFD}, val...)
	} else {
		return []byte{byte(n)}
	}
}

func getInputBinary(in models.TxInput) []byte {
	bin := make([]byte, 0)
	bin = append(bin, in.Hash...)

	index := make([]byte, 4)
	binary.LittleEndian.PutUint32(index, uint32(in.Index))
	bin = append(bin, index...)

	scriptLength := varint(uint64(len(in.Script)))
	bin = append(bin, scriptLength...)

	bin = append(bin, in.Script...)

	sequence := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequence, uint32(in.Sequence))
	bin = append(bin, sequence...)

	return bin
}

func getOutputBinary(out models.TxOutput) []byte {
	bin := make([]byte, 0)

	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, uint64(out.Value))
	bin = append(bin, value...)

	scriptLength := varint(uint64(len(out.Script)))
	bin = append(bin, scriptLength...)

	bin = append(bin, out.Script...)

	return bin
}

// Compute transaction hash
func putTransactionHash(tx *models.Transaction) {
	if tx.Hash != nil {
		return
	}

	bin := make([]byte, 0)

	//hasScriptWitness := tx.HasWitness()

	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(tx.NVersion))
	bin = append(bin, version...)

	//var flags byte
	//if hasScriptWitness {
	//	bin = append(bin, 0)
	//	flags |= 1
	//	bin = append(bin, flags)
	//}

	vinLength := varint(uint64(len(tx.Vin)))
	bin = append(bin, vinLength...)
	for _, in := range tx.Vin {
		bin = append(bin, getInputBinary(in)...)
	}

	voutLength := varint(uint64(len(tx.Vout)))
	bin = append(bin, voutLength...)
	for _, out := range tx.Vout {
		bin = append(bin, getOutputBinary(out)...)
	}

	//if hasScriptWitness {
	//	for _, in := range tx.Vin {
	//		bin = append(bin, in.ScriptWitnessBinary()...)
	//	}
	//}

	locktime := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktime, tx.Locktime)
	bin = append(bin, locktime...)

	tx.Hash = DoubleSha256(bin)
}

func (btc *Btc) parseBlockTransactionsFromFile(blockFile *parser.BlockFile) error {
  // Read transaction count to know how many transactions to parse
  TransactionCount := blockFile.ReadVarint()
  // TODO: TransactionCount : Compare with Index stored value (NTx)
  // log.Info(fmt.Sprintf("Total txns: %d vs %d", TransactionCount, btc.NTx))
  btc.Transactions = nil
  for t := uint64(0); t < TransactionCount; t++ {
    tx, err := btc.parseBlockTransactionFromFile(blockFile)
    if err != nil {
      return err
    }
    putTransactionHash(tx)
    log.Info(fmt.Sprintf("Transaction hash: %x", reverseHex(tx.Hash)))
    btc.Transactions = append(btc.Transactions, *tx)
  }
  return nil
}

func (btc *Btc) getTransaction() {
  f, err := db.GetFlag(btc.IndexDb, []byte("txindex"))
  if err != nil {
    log.Warn(err)
  }
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
