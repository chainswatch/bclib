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

/*
func txAbstractDb(filename string) (func(b *models.Block) error, error) {
	db, err := leveldb.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	db.NewIterator(util.BytesPrefix([]byte("n")), nil)
	lastKey := iter.Last()
	iter.Release()
	err = iter.Error()
	if err != nil {
	}

	return func(b *models.Block) error {
		batch := new(leveldb.Batch)
		for _, tx := range b.Txs {
			batch.Put("", tx.Hash)
			batch.Put([]byte("") + tx.Hash, "")
		}
		err = db.Write(batch, nil)
		return nil
	}, nil
}

// TODO: Move this function to application repository?
func csvExport(filename string) (func(b *models.Block) error, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	// defer file.Close()
	writer := csv.NewWriter(file)

	return func(b *models.Block) error {
		head := []string{fmt.Sprintf("%d", b.NHeight)}
		record := make([]string, 0)
		for _, tx := range b.Txs {
			for _, in := range tx.Vin {
				record = append(head, string(in.Hash))
			}
			record = append(record, string(tx.Hash))
			if err := writer.Write(record); err != nil {
				return err
			}
		}
		writer.Flush()
		return nil
	}, nil
}
*/

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

	indexDb, err := OpenIndexDb() // TODO: Error handling
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

	// Test loadFile
	db, err := OpenIndexDb()
	if err != nil {
		t.Fatal(err)
	}
	err = loadFile(db, 0, 10, dummyFunc, "output.csv")
	if err != nil {
		t.Fatal(err)
	}
}
