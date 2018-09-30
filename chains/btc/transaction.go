package btc

import (
  "app/misc"
  "app/models"
  "app/chains/parser"
  db "app/chains/repository"
  log "github.com/sirupsen/logrus"
  "encoding/binary"
  "fmt"
  "bytes"
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

type BlockReader interface {
  Peek(int) ([]byte, error)
  ReadByte() byte
  ReadBytes(uint64) []byte
  ReadVarint() uint64
  ReadUint32() uint32
  ReadUint64() uint64
  ReadInt32() int32
}

func parseBlockTransactionFromFile(br BlockReader) (*models.Transaction, error) {
  var err error = nil
  emptyByte := make([]byte, 32)
  allowWitness := true // TODO: Port code - !(s.GetVersion() & SERIALIZE_TRANSACTION_NO_WITNESS);

  tx := &models.Transaction{}
  // tx.StartPos = uint64(curPos)
  tx.NVersion = br.ReadInt32() // 

  // Check for extended transaction serialization format
  // var txInputLength uint64
  var txFlag byte
  // Try to read. Look for dummy
  p, _ := br.Peek(1)
  if p[0] == 0 {
    // We are dealing with extended transaction
    br.ReadByte()          // dummy (0x00)
    txFlag = br.ReadByte() // flags (!=0)
    if (txFlag != 0) {
      tx.NVin = uint32(br.ReadVarint())
    }
  } else {
    tx.NVin = uint32(br.ReadVarint())
  }

  for i := uint32(0); i < tx.NVin; i++ {
    input := models.TxInput{}
    input.Hash = br.ReadBytes(32) // 
    input.Index = br.ReadUint32() // TODO: Not sure if correctly read
    if input.Index == 0xFFFFFFFF && !bytes.Equal(input.Hash, emptyByte) { // block-reward case
      log.Fatal("If Index is 0xFFFFFFFF, then Hash should be nil. ",
      " Input: ", input.Index,
      " Hash: ", input.Hash)
    }
    scriptLength := br.ReadVarint()
    input.Script = br.ReadBytes(scriptLength)
    input.Sequence = br.ReadUint32()
    if (input.Sequence != 0xFFFFFFFF) { // Always 0xFFFFFFFF check
      log.Warn("Input Sequence != 0xFFFFFFFF: ", input.Sequence)
      err = fmt.Errorf("Input Sequence != 0xFFFFFFFF: %x", input.Sequence)
    }
    tx.Vin = append(tx.Vin, input)
  }

  tx.NVout = uint32(br.ReadVarint())
  for i := uint32(0); i < tx.NVout; i++ {
    output := models.TxOutput{}
    output.Index = i
    output.Value = int64(br.ReadUint64())
    scriptLength := br.ReadVarint()
    output.Script = br.ReadBytes(scriptLength)
    // output.hash160 // TODO
    getAddress(output.Script)
    tx.Vout = append(tx.Vout, output)
  }

  if (txFlag & 1) == 1 && allowWitness {
    txFlag ^= 1 // Not sure what this is for
    for i := uint32(0); i < tx.NVin; i++ {
      witnessCount := br.ReadVarint()
      tx.Vin[i].ScriptWitness = make([][]byte, witnessCount)
      for j := uint64(0); j < witnessCount; j++ {
        length := br.ReadVarint()
        tx.Vin[i].ScriptWitness[j] = br.ReadBytes(length)
      }
    }
  }

  tx.Locktime = br.ReadUint32() // Always 0 check
  if tx.Locktime != 0 {
    log.Fatal("Locktime is not 0")
  }

  return tx, err
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

	tx.Hash = misc.DoubleSha256(bin)
}

func (btc *Btc) parseBlockTransactionsFromFile(blockFile *parser.BlockFile) error {
  btc.Transactions = nil
  for t := uint32(0); t < btc.NTx; t++ {
    tx, err := parseBlockTransactionFromFile(blockFile)
    putTransactionHash(tx)
    if err != nil {
      log.Warn(fmt.Sprintf("txHash: %x", misc.ReverseHex(tx.Hash)))
    }
    tx.NVout = uint32(len(tx.Vout))
    // log.Info(fmt.Sprintf("Transaction hash: %x", misc.ReverseHex(tx.Hash)))
    // log.Fatal("Transaction")
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
    log.Warn("txindex is not enabled for your bitcoind")
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
