package btc

import (
	"git.posc.in/cw/bclib/serial"
	"git.posc.in/cw/bclib/models"
	"git.posc.in/cw/bclib/parser"

	log "github.com/sirupsen/logrus"
	"encoding/binary"
	"fmt"
)

const (
	blockMagicID = 0xd9b4bef9
	serGetHash = 1 << 2
)

func putBlockHash(b *models.BlockHeader) {
	bin := make([]byte, 0) // TODO: Optimize. 4 + 4 + 4 + 8 + 4 + 4

	value := make([]byte, 4)
	binary.LittleEndian.PutUint32(value, b.NVersion)
	bin = append(bin, value...)

	bin = append(bin, b.HashPrevBlock...)
	bin = append(bin, b.HashMerkleRoot...)

	binary.LittleEndian.PutUint32(value, b.NTime)
	bin = append(bin, value...)

	binary.LittleEndian.PutUint32(value, b.NBits)
	bin = append(bin, value...)

	binary.LittleEndian.PutUint32(value, b.NNonce)
	bin = append(bin, value...)

	b.HashBlock = serial.DoubleSha256(bin)
}

// TODO: Currently won't return any error
func decodeBlockHeader(bh *models.BlockHeader, br parser.Reader) {
	bh.NVersion = br.ReadUint32()
	bh.HashPrevBlock = br.ReadBytes(32) // FIXME: Slice out of bound (in production)
	bh.HashMerkleRoot = br.ReadBytes(32)
	bh.NTime = br.ReadUint32()
	bh.NBits = br.ReadUint32() // TODO: Parse this as mantissa?
	bh.NNonce = br.ReadUint32()
	putBlockHash(bh)
}

func decodeBlockHeaderIdx(br parser.Reader) (bh models.BlockHeader) {
  // Discard first varint
  // FIXME: Not exactly sure why need to, but if we don't do this we won't get correct values
	br.ReadVarint() // SerGetHash = 1 << 2

  bh.NHeight = uint32(br.ReadVarint())
  bh.NStatus = uint32(br.ReadVarint())
  bh.NTx = uint32(br.ReadVarint())
  if bh.NStatus & (blockHaveData|blockHaveUndo) > 0 {
    bh.NFile = uint32(br.ReadVarint())
  }
  if bh.NStatus & blockHaveData > 0 {
    bh.NDataPos = uint32(br.ReadVarint())
  }
  if bh.NStatus & blockHaveUndo > 0 {
    bh.NUndoPos = uint32(br.ReadVarint())
  }
	decodeBlockHeader(&bh, br)
	return
}

func decodeBlockTxs(b *models.Block, br parser.Reader) error {
	b.Txs = nil

	b.NTx = uint32(br.ReadCompactSize()) // TODO: Move outside of blockHeader?
	b.Txs = make([]models.Tx, b.NTx)
	for t := uint32(0); t < b.NTx; t++ {
		tx, err := DecodeTx(br)
		// putTxHash(tx) // TODO: Already done inside DecodeTx
		if err != nil {
			log.Warn(fmt.Sprintf("DecodeBlocksTxs(): txHash: %x", serial.ReverseHex(tx.Hash)))
			return err
		}
		tx.NVout = uint32(len(tx.Vout))
		b.Txs[t] = *tx
	}
	return nil
}

// DecodeBlock decodes a block
func DecodeBlock(br parser.Reader) (b *models.Block, err error) {
	b = &models.Block{}
	if br.Type() == "file" {
		magicID := uint32(br.ReadUint32())
		if magicID == 0 {
			return nil, fmt.Errorf("DecodeBlock: EOF")
		} else if magicID != blockMagicID {
			// blockFile.Seek(curPos, 0) // Restore pos before the error
			return nil, fmt.Errorf("Invalid block header: Can't find Magic ID")
		}
		b.NSize = br.ReadUint32() // Only for block files
	}

	decodeBlockHeader(&b.BlockHeader, br)
	err = decodeBlockTxs(b, br)
	if err != nil {
		return nil, err
	}

	if b.NHeight == 0 && len(b.Txs[0].Vin[0].Script) > 4 {
		cbase := b.Txs[0].Vin[0].Script[0:5]
		if cbase[0] == 3 {
			cbase[4] = 0
		}
		b.NHeight = binary.LittleEndian.Uint32(cbase[1:])
	}
	return b, err
}
