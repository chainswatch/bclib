package chains

import (
  "fmt"
)

func (btc *BtcBlock) ParseBlockTransactionsFromFile(blockFile *BlockFile) error {
	// Read transaction count to know how many transactions to parse
  TransactionCount := blockFile.ReadVarint()
  // TODO: TransactionCount : Compare with Index stored value (NTx)
	fmt.Printf("Total txns: %d vs %d\n", TransactionCount, btc.NTx)
	for t := uint64(0); t < TransactionCount; t++ {
		tx, err := btc.ParseBlockTransactionFromFile(blockFile)
		if err != nil {
			return err
		}
		btc.Transactions = append(btc.Transactions, *tx)
	}

	return nil
}

func (btc *BtcBlock) ParseBlockTransactionFromFile(blockFile *BlockFile) (*Transaction, error) {
	curPos, err := blockFile.Seek(0, 1)
	if err != nil {
		return nil, err
	}

	allowWitness := true // TODO: Port code - !(s.GetVersion() & SERIALIZE_TRANSACTION_NO_WITNESS);

	tx := &Transaction{}
	tx.StartPos = uint64(curPos)
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
