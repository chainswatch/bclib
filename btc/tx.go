package btc

import (
	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/parser"
	"github.com/chainswatch/bclib/serial"

	"bytes"
	"encoding/binary"
	"fmt"
)

// Witness : https://github.com/bitcoin/bitcoin/blob/master/src/primitives/transaction.h
// const serializeTransactionNoWitness = 0x40000000;

// DecodeTx decodes a transaction
func DecodeTx(br parser.Reader) (*models.Tx, error) {
	var txFlag byte // Check for extended transaction serialization format
	emptyByte := make([]byte, 32)
	allowWitness := true // TODO: Port code - !(s.GetVersion() & SERIALIZE_TRANSACTION_NO_WITNESS);
	tx := &models.Tx{}

	tx.NVersion = br.ReadInt32()
	tx.NVin = uint32(br.ReadCompactSize())
	if tx.NVin == 0 { // We are dealing with extended transaction (witness format)
		txFlag, _ = br.ReadByte() // TODO: Error handling
		if txFlag != 0x01 {       // Must be 1, other flags may be supported in the future
			return nil, fmt.Errorf("Witness tx but flag is %x != 0x01", txFlag)
		}
		tx.NVin = uint32(br.ReadCompactSize())
	}

	tx.Vin = make([]models.TxInput, tx.NVin)
	for i := uint32(0); i < tx.NVin; i++ {
		input := models.TxInput{}
		input.Hash = br.ReadBytes(32)                                         // Transaction hash in a prev transaction
		input.Index = br.ReadUint32()                                         // Transaction index in a prev tx TODO: Not sure if correctly read
		if input.Index == 0xFFFFFFFF && !bytes.Equal(input.Hash, emptyByte) { // block-reward case
			return nil, fmt.Errorf("If Index is 0xFFFFFFFF, then Hash should be nil. Input: %d, Hash: %x", input.Index, input.Hash)
		}
		scriptLength := br.ReadCompactSize()
		input.Script = br.ReadBytes(scriptLength)
		input.Sequence = br.ReadUint32()
		tx.Vin[i] = input
	}

	tx.NVout = uint32(br.ReadCompactSize())
	tx.Vout = make([]models.TxOutput, tx.NVout)
	for i := uint32(0); i < tx.NVout; i++ {
		output := models.TxOutput{}
		output.Index = i
		output.Value = br.ReadUint64()
		scriptLength := br.ReadCompactSize()
		output.Script = br.ReadBytes(scriptLength)
		output.AddrType, output.Addr = getPkeyFromScript(output.Script) // TODO: Optimize
		tx.Vout[i] = output
	}

	if (txFlag&1) == 1 && allowWitness {
		// txFlag ^= 1 // Not sure what this is for
		tx.Segwit = true
		for i := uint32(0); i < tx.NVin; i++ {
			witnessCount := br.ReadCompactSize()
			tx.Vin[i].ScriptWitness = make([][]byte, witnessCount)
			for j := uint64(0); j < witnessCount; j++ {
				length := br.ReadCompactSize()
				tx.Vin[i].ScriptWitness[j] = br.ReadBytes(length)
			}
		}
	} // TODO: Missing 0 field?

	tx.Locktime = br.ReadUint32()
	putTxHash(tx)
	return tx, nil
}

func getInputBinary(in models.TxInput) []byte {
	bin := make([]byte, 0)
	bin = append(bin, in.Hash...)

	index := make([]byte, 4)
	binary.LittleEndian.PutUint32(index, uint32(in.Index))
	bin = append(bin, index...)

	scriptLength := parser.CompactSize(uint64(len(in.Script)))
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

	scriptLength := parser.CompactSize(uint64(len(out.Script)))
	bin = append(bin, scriptLength...)

	bin = append(bin, out.Script...)

	return bin
}

// 0100000001e507cb947464fc74540a9c197f815aa283ba9db74185ac08449c38491a8c34ac00000000
// Compute transaction hash ( [nVersion][Inputs][Outputs][nLockTime] )
func putTxHash(tx *models.Tx) {
	bin := make([]byte, 0)
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(tx.NVersion))
	bin = append(bin, version...)

	vinLength := parser.CompactSize(uint64(tx.NVin))
	bin = append(bin, vinLength...)
	for _, in := range tx.Vin {
		bin = append(bin, getInputBinary(in)...)
	}

	voutLength := parser.CompactSize(uint64(tx.NVout))
	bin = append(bin, voutLength...)
	for _, out := range tx.Vout {
		bin = append(bin, getOutputBinary(out)...)
	}

	locktime := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktime, tx.Locktime)
	bin = append(bin, locktime...)

	tx.Hash = serial.DoubleSha256(bin)
}
