package btc

import (
  "git.posc.in/cw/watchers/serial"
  "git.posc.in/cw/watchers/models"
  "git.posc.in/cw/watchers/parser"
  log "github.com/sirupsen/logrus"
  "encoding/binary"
  "fmt"
  "bytes"
)

// const serializeTransactionNoWitness = 0x40000000; // Witness : https://github.com/bitcoin/bitcoin/blob/master/src/primitives/transaction.h

/*
* Basic transaction serialization format:
* - int32_t nVersion          / Expected to be 1, but can be garbage in rare cases
* - std::vector<CTxIn> vin
* - std::vector<CTxOut> vout
* - uint32_t nLockTime        / Always set to 0
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

// DecodeTx decode a transaction
func DecodeTx(br parser.Reader) (*models.Transaction, error) {
  var err error
  var txFlag byte // Check for extended transaction serialization format
  emptyByte := make([]byte, 32)
  allowWitness := true // TODO: Port code - !(s.GetVersion() & SERIALIZE_TRANSACTION_NO_WITNESS);
  tx := &models.Transaction{}

  tx.NVersion = br.ReadInt32()
  log.Debug("NVersion:", tx.NVersion)
  tx.NVin = uint32(br.ReadVarint())
  if tx.NVin == 0 { // We are dealing with extended transaction (witness format)
    log.Debug("Segwit Transaction")
    txFlag = br.ReadByte()
    log.Debug("Flag: ", txFlag)
    if (txFlag != 0x01) { // Must be 1, other flags may be supported in the future
      log.Fatal("Witness tx but flag is ", txFlag, " != 0x01")
    }
    tx.NVin = uint32(br.ReadVarint())
  }

  log.Debug("Number of inputs:", tx.NVin)
  for i := uint32(0); i < tx.NVin; i++ {
    input := models.TxInput{}
    input.Hash = br.ReadBytes(32) // Transaction hash in a prev transaction
    input.Index = br.ReadUint32() // Transaction index in a prev tx TODO: Not sure if correctly read
    if input.Index == 0xFFFFFFFF && !bytes.Equal(input.Hash, emptyByte) { // block-reward case
      log.Fatal("If Index is 0xFFFFFFFF, then Hash should be nil. ",
      " Input: ", input.Index,
      " Hash: ", input.Hash)
    }
    scriptLength := br.ReadVarint()
    input.Script = br.ReadBytes(scriptLength)
    input.Sequence = br.ReadUint32()
    tx.Vin = append(tx.Vin, input)
  }

  tx.NVout = uint32(br.ReadVarint())
  log.Debug("Number of outputs:", tx.NVout)
  for i := uint32(0); i < tx.NVout; i++ {
    output := models.TxOutput{}
    output.Index = i
    output.Value = int64(br.ReadUint64())
    log.Debug("Value ", i, ": ", output.Value)
    scriptLength := br.ReadVarint()
    output.Script = br.ReadBytes(scriptLength)
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
  } // TODO: Missing 0 field?

  tx.Locktime = br.ReadUint32()
  log.Debug("Locktime: ", tx.Locktime)
	putTransactionHash(tx)
  if err != nil {
		log.Info(fmt.Sprintf("txHash: %x", serial.ReverseHex(tx.Hash)))
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

/*
func hasWitness() {
   tx.Vin[i].ScriptWitness = make([][]byte, witnessCount)
}
*/

// Compute transaction hash
func putTransactionHash(tx *models.Transaction) {
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

	tx.Hash = serial.DoubleSha256(bin)
}

// Reads a transaction directly from a file
func (btc *btc) parseBlockTransactionsFromFile(blockFile *parser.File) error {
  btc.Transactions = nil
  for t := uint32(0); t < btc.NTx; t++ {
    tx, err := DecodeTx(blockFile)
    putTransactionHash(tx)
    if err != nil {
      log.Warn(fmt.Sprintf("txHash: %x", serial.ReverseHex(tx.Hash)))
    }
    tx.NVout = uint32(len(tx.Vout))
    btc.Transactions = append(btc.Transactions, *tx)
  }
  return nil
}
