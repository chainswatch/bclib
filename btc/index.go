package btc

import (
	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/parser"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"encoding/binary"
	"os"
)

type blockFile struct {
	fileNum     uint32
	NBlocks     uint32
	Size        uint32
	UndoSize    uint32
	HeightFirst uint32
	HeightLast  uint32
	TimeFirst   uint32
	TimeLast    uint32
}

func decodeBlockFileIdx(br parser.Reader) *blockFile {
	f := &blockFile{}
	f.NBlocks = uint32(br.ReadVarint())
	f.Size = uint32(br.ReadVarint())
	f.UndoSize = uint32(br.ReadVarint())
	f.HeightFirst = uint32(br.ReadVarint())
	f.HeightLast = uint32(br.ReadVarint())
	f.TimeFirst = uint32(br.ReadVarint())
	f.TimeLast = uint32(br.ReadVarint())
	return f
}

// Get block index record by hash
func blockIndexRecord(db *leveldb.DB, h []byte) (bh *models.BlockHeader, err error) {
	data, err := db.Get(append([]byte("b"), h...), nil)
	if err != nil {
		return
	}
	buf, err := parser.New(data)
	if err != nil {
		return
	}
	return decodeBlockHeaderIdx(buf), err
}

// Get file index record by number
func fileIndexRecord(db *leveldb.DB, n uint32) (*blockFile, error) {
	value := make([]byte, 4)
	binary.BigEndian.PutUint32(value, n)
	data, err := db.Get(append([]byte("f"), value...), nil)
	if err != nil {
		return nil, err
	}
	buf, err := parser.New(data)
	if err != nil {
		return nil, err
	}
	f := decodeBlockFileIdx(buf)
	return f, nil
}

// Get file information record by number

// Get last block file number used

// GetFlag checks wether txindex is enabled or not
func GetFlag(db *leveldb.DB, name []byte) (bool, error) {
	command := append([]byte("F"), byte(len(name)))
	command = append(command, name...)
	data, err := db.Get(command, nil)
	if err != nil {
		return false, err
	}

	return data[0] == []byte("1")[0], nil
}

// decode block header from index files
func decodeBlockHeaderIdx(br parser.Reader) *models.BlockHeader {
	bh := new(models.BlockHeader)

	br.ReadVarint() // SerGetHash = 1 << 2 (client version)

	bh.NHeight = uint32(br.ReadVarint())
	bh.NStatus = uint32(br.ReadVarint())
	bh.NTx = uint32(br.ReadVarint())
	if bh.NStatus&(blockHaveData|blockHaveUndo) == 0 {
		return nil
	}
	bh.NFile = uint32(br.ReadVarint())
	if bh.NStatus&blockHaveData > 0 {
		bh.NDataPos = uint32(br.ReadVarint())
	}
	if bh.NStatus&blockHaveUndo > 0 {
		bh.NUndoPos = uint32(br.ReadVarint())
	}

	decodeBlockHeader(bh, br)
	return bh
}

// OpenIndexDb gets transaction index record
func OpenIndexDb() (*leveldb.DB, error) {
	dataDir := os.Getenv("DATADIR")
	db, err := leveldb.OpenFile(dataDir+"/blocks/index", &opt.Options{
		ReadOnly: true,
	})
	return db, err
}
