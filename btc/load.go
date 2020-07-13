package btc

import (
	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/parser"

	"fmt"
	log "github.com/sirupsen/logrus"

	"os"
	"sort"
	"strconv"
	"strings"
)

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
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		if k > oldb.NFile || k+1 >= bh.NFile {
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

func loadRawBlock() (func(height uint32) (parser.Reader, error), error) {
	lookup, err := LoadHeaderIndex()
	log.Info("Index is built: ", len(lookup))
	if err != nil {
		return nil, err
	}
	files := make(map[uint32]parser.Reader) // map[BlockHeight]

	return func(height uint32) (parser.Reader, error) {
		bh, exist := lookup[height]
		if !exist {
			return nil, fmt.Errorf("LoadBlock(): File for height %d does not exist", height)
		}
		if bh.NHeight != height {
			return nil, fmt.Errorf("LoadBlock(): Loaded header has wrong height %d != %d", bh.NHeight, height)
		}
		file, exist := files[bh.NFile]
		if !exist { // file open ?
			log.Info(fmt.Sprintf("Height: %d File: %d Length(files)= %d", bh.NHeight, bh.NFile, len(files)))
			buf, err := parser.New(bh.NFile)
			if err != nil {
				return nil, err
			}
			files[bh.NFile] = buf
			file = buf
			if err = closeOldFile(bh, lookup, files); err != nil {
				return nil, err
			}
		}
		file.Seek(int64(bh.NDataPos-8), 0)

		return file, nil
	}, nil
}

// LoadBlockToFile prints block content in a file
func LoadBlockToFile(path string, height uint32) error {
	load, err := loadRawBlock()
	if err != nil {
		return err
	}
	file, err := load(height)
	file.Seek(4, 0)
	nSize := file.ReadUint32()
	log.Info("Size: ", nSize) // Only for block files
	content := file.ReadBytes(uint64(nSize))

	fout, err := os.Create(fmt.Sprintf("%s/block%d.dat", path, height))
	if err != nil {
		return err
	}
	if _, err := fout.Write(content); err != nil {
		return err
	}
	return fout.Close()
}

func LoadBlock() (func(height uint32) (*models.Block, error), error) {
	load, err := loadRawBlock()
	if err != nil {
		return nil, err
	}
	return func(height uint32) (*models.Block, error) {
		file, err := load(height)
		if err != nil {
			return nil, err
		}

		b, err := DecodeBlock(file)
		if err != nil {
			return nil, fmt.Errorf("LoadBlock(): Height %d: %s", height, err.Error())
		}
		b.NHeight = height // FIXME: DecodeBlock does not work for genesis block

		return b, nil
	}, nil
}

type apply func(interface{}) (func(b *models.Block) error, error)

// LoadFile allows to traverse the blocks by height order while applying a function argFn
func LoadFile(fromh, toh uint32, newFn apply, argFn interface{}) error {
	loadBlock, err := LoadBlock()
	if err != nil {
		return err
	}
	files := make(map[uint32]parser.Reader) // map[BlockHeight]

	fn, err := newFn(argFn)
	if err != nil {
		return err
	}

	for h := fromh; h <= toh; h++ {
		b, err := loadBlock(h)
		if err != nil {
			return err
		}

		if err = fn(b); err != nil {
			if strings.HasPrefix(err.Error(), "Jump to height ") {
				s := strings.TrimPrefix(err.Error(), "Jump to height ")
				tmp, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return err
				}
				h = uint32(tmp) - 1
				continue
			}
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
