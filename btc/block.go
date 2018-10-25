package btc

import (
  "git.posc.in/cw/watchers/serial"
  "git.posc.in/cw/watchers/models"
  "git.posc.in/cw/watchers/parser"

  log "github.com/sirupsen/logrus"
  "time"
  "fmt"
)

// Parse the header fields except the MagicId
// TODO: Currently won't return any error
func decodeBlockHeader(btc *models.Block, br parser.Reader) {
  // btc.Length = br.ReadUint32() // Maybe only for raw files?

  btc.NVersion = br.ReadUint32()
  btc.HashPrevBlock = br.ReadBytes(32)
  btc.HashMerkleRoot = br.ReadBytes(32)
  btc.NTime = time.Unix(int64(br.ReadUint32()), 0)
  btc.NBits = br.ReadUint32() // TODO: Parse this as mantissa?
  btc.NNonce = br.ReadUint32()
  btc.NTx = uint32(br.ReadVarint())
}

func decodeBlockTxs(btc *models.Block, br parser.Reader) error {
  btc.Txs = nil
  for t := uint32(0); t < btc.NTx; t++ {
    tx, err := DecodeTx(br)
    putTransactionHash(tx)
    if err != nil {
      log.Warn(fmt.Sprintf("DecodeBlocksTxs(): txHash: %x", serial.ReverseHex(tx.Hash)))
      return err
    }
    tx.NVout = uint32(len(tx.Vout))
    btc.Txs = append(btc.Txs, *tx)
  }
  return nil
}

// DecodeBlock decodes a block
func DecodeBlock(br parser.Reader) (*models.Block, error) {
  btc := &models.Block{}
  decodeBlockHeader(btc, br)
  err := decodeBlockTxs(btc, br)
  return btc, err
}
