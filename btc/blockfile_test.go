package btc

import (
	"git.posc.in/cw/bclib/models"
	"git.posc.in/cw/bclib/parser"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/joho/godotenv"
	"testing"
	"os"

	_ "net/http/pprof"
	"net/http"
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

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()


	indexDb, err := OpenIndexDb()
	if err != nil {
		t.Fatal(err)
	}

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
	indexDb.Close()

	// Open blockfile 0 and parse blocks
	buf, err := parser.New(uint32(0))
	if err != nil {
		t.Fatal(err)
	}
	buf.Reset()

	b := &models.Block{}
	i := 0
	for { // TODO: Test EOF
		if err := DecodeBlock(b, buf); err != nil {
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

	/*
	err = LoadFile(0, 100000, dummyFunc, "")
	if err != nil {
		t.Fatal(err)
	}

	err = LoadFile(200000, 201000, dummyFunc, "")
	if err != nil {
		t.Fatal(err)
	}
	*/
}
