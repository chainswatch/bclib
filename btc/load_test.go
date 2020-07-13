package btc

import (
	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/parser"

	"github.com/joho/godotenv"

	"os"
	"testing"

	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func dummyFunc(_ interface{}) (func(b *models.Block) error, error) {
	return func(b *models.Block) error {
		return nil
	}, nil
}

func jumpFunc(_ interface{}) (func(b *models.Block) error, error) {
	var min uint32 = 201000
	return func(b *models.Block) error {
		if b == nil {
			return nil
		}
		if b.NHeight < min {
			return fmt.Errorf("Jump to height %d", min)
		}
		return nil
	}, nil
}

func TestLoad(t *testing.T) {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := LoadBlockToFile(".", 265458); err != nil {
		t.Error(err)
	}

	indexDb, err := OpenIndexDb()
	if err != nil {
		t.Fatal(err)
	}

	// Iter over leveldb block index
	iter := indexDb.NewIterator(util.BytesPrefix([]byte("b")), nil)
	for i := 0; i < 200; i++ {
		iter.Next()
		// hashBlock := iter.Key()
		data := iter.Value()
		buf, err := parser.New(data)
		if err != nil {
			t.Fatal(err)
		}
		decodeBlockHeaderIdx(buf)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Warn(err)
	}
	indexDb.Close()

	// Open blockfile 0 and parse blocks
	buf, err := parser.New(uint32(0))
	if err != nil {
		t.Fatal(err)
	}
	buf.Reset()

	i := 0
	for { // TODO: Test EOF
		_, err := DecodeBlock(buf)
		if err != nil {
			if err.Error() != "DecodeBlock: EOF" {
				t.Error(err)
			}
			break
		}
		i++
	}
	if i < 110000 {
		t.Errorf("Only %d blocks read in blockfile 0", i)
	}

	loadBlock, err := LoadBlock()
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = loadBlock(0); err != nil {
		t.Fatal(err)
	}
	if _, _, err = loadBlock(442019); err != nil {
		t.Fatal(err)
	}

	err = LoadFile(0, 100000, dummyFunc, "")
	if err != nil {
		t.Fatal(err)
	}

	err = LoadFile(0, 110000, dummyFunc, "")
	if err != nil {
		t.Fatal(err)
	}

	err = LoadFile(100000, 120000, dummyFunc, "")
	if err != nil {
		t.Fatal(err)
	}

	err = LoadFile(200000, 202000, jumpFunc, "")
	if err != nil {
		t.Fatal(err)
	}

	// Used to be an error (because of orphan blocks)
	err = LoadFile(478500, 478700, dummyFunc, "")
	if err != nil {
		t.Fatal(err)
	}
}
