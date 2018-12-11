package btc

import (
	"git.posc.in/cw/bclib/models"
	"git.posc.in/cw/bclib/parser"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/joho/godotenv"
	"testing"
	"os"
)

func dummyFunc(_ string) (func(b *models.Block) error, error) {
	return func(b *models.Block) error {
		return nil
	}, nil
}

func TestBlockFile(t *testing.T) {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			t.Fatal(err)
		}
	}

	indexDb, err := OpenIndexDb()
	if err != nil {
		t.Fatal(err)
	}
	defer indexDb.Close()

	// Iter over leveldb block index
	iter := indexDb.NewIterator(util.BytesPrefix([]byte("b")), nil)
	for i := 0; i < 100; i++ {
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

	// Open blockfile 0 and parse blocks
	buf, err := parser.New(uint32(0))
	if err != nil {
		t.Fatal(err)
	}
	buf.Reset()

	b := &models.Block{}
	for i := 0; i <= 100000; i++ { // TODO: Test EOF
		err := DecodeBlock(b, buf)
		if err != nil {
			t.Error(err)
		}
	}

	err = LoadFile(indexDb, 0, 100000, dummyFunc, "")
	if err != nil {
		t.Fatal(err)
	}

	err = LoadFile(indexDb, 200000, 201000, dummyFunc, "")
	if err != nil {
		t.Fatal(err)
	}
}
