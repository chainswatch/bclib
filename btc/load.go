package btc

import (
	"git.posc.in/cw/bclib/models"
	"git.posc.in/cw/bclib/parser"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb/util"
	"fmt"

	"sort"
)

// constructs a map of the form map[BlockHeight] = BlockHeader.
// In particular, BlockHeader contains DataPos and FileNum
func loadHeaderIndex() (lookup map[uint32]*models.BlockHeader, err error) {
	db, err := OpenIndexDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	iter := db.NewIterator(util.BytesPrefix([]byte("b")), nil)
	lookup = make(map[uint32]*models.BlockHeader)
	for iter.Next() {
		// hashBlock := iter.Key()
		data := iter.Value()
		buf, err := parser.New(data)
		if err != nil {
			return nil, err
		}
		tmp := decodeBlockHeaderIdx(buf)
		lookup[tmp.NHeight] = &tmp
	}
	iter.Release()
	return lookup, nil
}

type apply func(string) (func(b *models.Block) error, error)

func closeOldFile(bh *models.BlockHeader, lookup map[uint32]*models.BlockHeader, files map[uint32]parser.Reader) error {
	if bh.NHeight < 2048 {
		return nil
	}
	oldh := bh.NHeight - 2048
	oldb, exist := lookup[oldh]
	if !exist {
		return fmt.Errorf("closeOldFile: Could not find old file for height %d", oldh)
	}
	var keys []uint32
	for k := range files {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {return keys[i] < keys[j]})
	for _, k := range keys {
		if k > oldb.NFile || k + 1 >= bh.NFile {
			break
		}
		oldf, exist := files[k]
		if !exist {
			return fmt.Errorf("closeOldFile: Could not find old file reader. Height %d File %d", oldh, k) 
		}
		oldf.Close()
		delete(files, k)
	}
	return nil
}

// LoadFile allows to traverse the blocks by height order while applying a function argFn
func LoadFile(fromh, toh uint32, newFn apply, argFn string) error {
	lookup, err := loadHeaderIndex()
	log.Info("Index is built: ", len(lookup))
	if err != nil {
		return err
	}
	files := make(map[uint32]parser.Reader) // map[BlockHeight]

	fn, err := newFn(argFn)
	if err != nil {
		return err
	}

	var bh = &models.BlockHeader{}
	var exist bool
	for h := fromh; h <= toh; h++ {
		bh, exist = lookup[h]
		if !exist { // header ?
			return fmt.Errorf("File for height %d does not exist", h)
		}
		if bh.NHeight != h {
			return fmt.Errorf("Loaded header has wrong height %d != %d", bh.NHeight, h)
		}
		file, exist := files[bh.NFile]
		if !exist { // file open ?
			log.Info(fmt.Sprintf("Height: %d File: %d Length(files)= %d", bh.NHeight, bh.NFile, len(files)))
			buf, err := parser.New(bh.NFile)
			if err != nil {
				return err
			}
			files[bh.NFile] = buf
			file = buf
			if err = closeOldFile(bh, lookup, files); err != nil {
				return err
			}
		}
		file.Seek(int64(bh.NDataPos - 8), 0)
		b, err := DecodeBlock(file)
		if err != nil {
			return fmt.Errorf("File %d, height %d: %s", b.NFile, h, err.Error())
		}
		b.NHeight = h // FIXME: DecodeBlock does not work for genesis block
		if err = fn(b); err != nil {
			return err
		}
		// TODO: Check number of file open (always <= 2)
	}
	log.Info("Number of files still open: ", len(files))
	for _, value := range files {
		value.Close()
	}
	return fn(nil) // Signal fn to close
}
