package btc

import (
  "git.posc.in/cw/watchers/parser"
  "git.posc.in/cw/watchers/serial"

  log "github.com/sirupsen/logrus"
  "time"
  "fmt"
)

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func decodeBlockHeader(btc *btc, br parser.Reader) {
  btc.Length = br.ReadUint32()
  btc.NVersion = br.ReadInt32() // TODO: Uint32? (Always 1)
  btc.HashPrevBlock = br.ReadBytes(32)
  btc.HashMerkleRoot = br.ReadBytes(32)
  btc.NTime = time.Unix(int64(br.ReadUint32()), 0)
  btc.TargetDifficulty = br.ReadUint32() // TODO: Parse this as mantissa?
  btc.NNonce = br.ReadUint32()
  btc.NTx = uint32(br.ReadVarint())
}

func decodeBlockTxs(btc *btc, br parser.Reader) error {
  btc.Transactions = nil
  for t := uint32(0); t < btc.NTx; t++ {
    tx, err := DecodeTx(br)
    putTransactionHash(tx)
    if err != nil {
      log.Warn(fmt.Sprintf("txHash: %x", serial.ReverseHex(tx.Hash)))
    }
    tx.NVout = uint32(len(tx.Vout))
    btc.Transactions = append(btc.Transactions, *tx)
  }
  return nil
}

func DecodeBlock(br parser.Reader) {
  btcBlock := &btc{}
  decodeBlockHeader(btcBlock, br)

  err := decodeBlockTxs(btcBlock, br)
  if err != nil {
    // blockFile.Seek(curPos, 0) // Seek back to original pos before we encounter the error
    log.Warn(err)
    // return
  }
}
