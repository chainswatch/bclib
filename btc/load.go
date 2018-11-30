package btc

import (
	"git.posc.in/cw/bclib/models"
	"git.posc.in/cw/bclib/parser"

	log "github.com/sirupsen/logrus"
  "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"fmt"
)

// constructs a map of the form map[BlockHeight] = BlockHeader.
// In particular, BlockHeader contains DataPos and FileNum
func loadHeaderIndex(db *leveldb.DB) (map[uint32]*models.Block, error) {
	iter := db.NewIterator(util.BytesPrefix([]byte("b")), nil)
	lookup := make(map[uint32]*models.Block)
	for iter.Next() {
		// hashBlock := iter.Key()
		data := iter.Value()
		buf, err := parser.New(data)
		if err != nil {
			return nil, err
		}
		tmp := decodeBlockHeaderIdx(buf)
		lookup[tmp.NHeight] = tmp
	}
	iter.Release()
	return lookup, nil
}

type apply func(string) (func(b *models.Block) error, error)

// TODO: This function should take a pointer to function as input
// Its main purpose is only to parse the blocks fromh toh.
func loadFile(db *leveldb.DB, fromh, toh uint32, newFn apply, argFn string) error {
	lookup, err := loadHeaderIndex(db)
	log.Info("Index is built: ", len(lookup))
	if err != nil {
		return err
	}
	files := make(map[uint32]parser.Reader) // map[BlockHeight]

	fn, err := newFn(argFn)
	if err != nil {
		return err
	}

	var b = &models.Block{}
	var exist bool
	for h := fromh; h <= toh; h++ {
		b, exist = lookup[h]
		if !exist { // header ?
			return fmt.Errorf("File for height %d does not exist", h)
		}
		if b.NHeight != h {
			return fmt.Errorf("Loaded header has wrong height %d != %d", b.NHeight, h)
		}
		file, exist := files[h]
		if !exist { // file open ?
			buf, err := parser.New(b.NFile)
			if err != nil {
				return err
			}
			files[h] = buf
			file = buf
		}
		file.Seek(int64(b.NDataPos - 8), 0)
		err = DecodeBlock(b, file)
		if err != nil {
			return fmt.Errorf("File %d, height %d: %s", b.NFile, h, err.Error())
		}
		if err = fn(b); err != nil {
			return err
		}
		// TODO: Close file if necessary
		// TODO: Check number of file open (always <= 2)
	}
	log.Info(b.NHeight)
	return nil
}
