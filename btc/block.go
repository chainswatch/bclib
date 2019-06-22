package btc

import (
	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/parser"
	"github.com/chainswatch/bclib/serial"

	"encoding/binary"
	"fmt"
)

const (
	blockMagicID = 0xd9b4bef9
	serGetHash   = 1 << 2
)

func putBlockHash(b *models.BlockHeader) {
	bin := make([]byte, 0) // TODO: Optimize. 4 + 4 + 4 + 8 + 4 + 4

	value := make([]byte, 4)
	binary.LittleEndian.PutUint32(value, b.NVersion) // 4
	bin = append(bin, value...)

	bin = append(bin, b.HashPrev...)       // ?
	bin = append(bin, b.HashMerkleRoot...) // ?

	binary.LittleEndian.PutUint32(value, b.NTime) // 4
	bin = append(bin, value...)

	binary.LittleEndian.PutUint32(value, b.NBits) // 4
	bin = append(bin, value...)

	binary.LittleEndian.PutUint32(value, b.NNonce) // 4
	bin = append(bin, value...)

	b.Hash = serial.DoubleSha256(bin)
}

// TODO: Currently won't return any error
func decodeBlockHeader(bh *models.BlockHeader, br parser.Reader) {
	bh.NVersion = br.ReadUint32()
	bh.HashPrev = br.ReadBytes(32) // FIXME: Slice out of bound (in production)
	bh.HashMerkleRoot = br.ReadBytes(32)
	bh.NTime = br.ReadUint32()
	bh.NBits = br.ReadUint32() // TODO: Parse this as mantissa?
	bh.NNonce = br.ReadUint32()
	putBlockHash(bh)
}

func decodeBlockTxs(b *models.Block, br parser.Reader) error {
	b.Txs = nil

	b.NTx = uint32(br.ReadCompactSize()) // TODO: Move outside of blockHeader?
	b.Txs = make([]models.Tx, b.NTx)
	for t := uint32(0); t < b.NTx; t++ {
		tx, err := DecodeTx(br)
		if err != nil {
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
	decodeBlockTxs(b, br)

	if b.NHeight == 0 && len(b.Txs[0].Vin[0].Script) > 4 {
		cbase := b.Txs[0].Vin[0].Script[0:5]
		if cbase[0] == 3 {
			cbase[4] = 0
		}
		b.NHeight = binary.LittleEndian.Uint32(cbase[1:])
	}
	return b, err
}
