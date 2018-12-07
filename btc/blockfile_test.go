package btc

import (
	"git.posc.in/cw/bclib/models"
	"git.posc.in/cw/bclib/parser"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
  "github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/joho/godotenv"
	"encoding/binary"
	"testing"
	"fmt"
	"os"
)

// TODO: Move this function to application repository?
/*
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

/*
func dummyFunc(_ string) (func(b *models.Block) error, error) {
	return func(b *models.Block) error {
		return nil
	}, nil
}
*/

// Create an abstract database of the transactions
func txAbstractDB(filename string) (func(b *models.Block) error, error) {
	var idx uint32 = 1 // 0 is reserved

	db, err := leveldb.OpenFile(filename, nil)
	if err != nil {
		return nil, err
	}
	iter := db.NewIterator(util.BytesPrefix([]byte("n")), nil)
	for iter.Next() {
		log.Info(fmt.Sprintf("%x", iter.Key()))
	}
	if iter.Last() { // load last idx if exists
		idx = binary.BigEndian.Uint32(iter.Key()[1:])
	} else {
		log.Info("Start from new")
		tkey := make([]byte, 1 + 32 + 4)
		tval := make([]byte, 4)
		tkey[0] = byte('t')
		binary.BigEndian.PutUint32(tval, 0)
		if err = db.Put(tkey, tval, nil); err != nil {
			log.Warn("Initial insert: ", err)
		}
	}
	log.Info("Start from idx = ", idx)
	iter.Release()
	if err = iter.Error(); err != nil {
		return nil, err
	}

	return func(b *models.Block) error {
		if b == nil {
			return db.Close()
		}
		nkey := make([]byte, 1 + 4)
		nval := make([]byte, 32 + 4 + 4) // hash + block height + NVin
		tkey := make([]byte, 1 + 32 + 4)
		tval := make([]byte, 4)
		nkey[0] = byte('n')
		tkey[0] = byte('t')

		batch := new(leveldb.Batch)
		for _, tx := range b.Txs {
			binary.BigEndian.PutUint32(nkey[1:], b.NHeight) // TODO: Opti
			copy(tkey[1:], tx.Hash)
			copy(tkey[33:], nkey[1:])
			binary.BigEndian.PutUint32(tval, idx)
			batch.Put(tkey, tval)

			// 'n' + id -> tx hash + block height + NVin + {id, id, id...}
			// height
			// tx.NVin
			// tx.Vin[].Hash
			copy(nval, tx.Hash)
			binary.BigEndian.PutUint32(nval[32:], b.NHeight) // TODO: Opti
			binary.BigEndian.PutUint32(nval[36:], tx.NVin)
			var i uint32
			for i = 0; i < tx.NVin; i++ {
				copy(tkey[1:], tx.Vin[i].Hash)
				iter := db.NewIterator(util.BytesPrefix(tkey[:33]), nil)
				if iter.Next() {
					index := make([]byte, 4)
					binary.BigEndian.PutUint32(index, tx.Vin[i].Index)
					data := iter.Value()
					nval = append(nval, data...)
					nval = append(nval, index...)
				} else { // TODO: Error
					log.Info(fmt.Sprintf("Key not found: %x %x", tkey[:33], tx.Vin[i].Hash))
				}
				iter.Release()
			}
			batch.Put(nkey, nval)
			idx++
		}
		return db.Write(batch, nil)
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

	// Test LoadFile with a dummyFunc
	db, err := OpenIndexDb()
	if err != nil {
		t.Fatal(err)
	}

	// Test on tmp storage
	err = LoadFile(db, 0, 1e1, txAbstractDB, "/tmp/abstracts")
	if err != nil {
		t.Fatal(err)
	}

  dbtx, err := leveldb.OpenFile("/tmp/abstracts", &opt.Options{
    ReadOnly: true,
  })
	if err != nil {
		t.Fatal(err)
	}
	iter = dbtx.NewIterator(util.BytesPrefix([]byte("n")), nil)
	ncount := 0
	for iter.Next() {
		log.Info(fmt.Sprintf("%x %x", iter.Key(), iter.Value()))
		ncount++
	}
	iter.Release()
	if err = iter.Error(); err != nil {
		t.Fatal(err)
	}

	iter = dbtx.NewIterator(util.BytesPrefix([]byte("t")), nil)
	tcount := 0
	for iter.Next() {
		log.Info(fmt.Sprintf("%x %x", iter.Key(), iter.Value()))
		tcount++
	}
	iter.Release()
	if err = iter.Error(); err != nil {
		t.Fatal(err)
	}
	if ncount != tcount {
		t.Fatalf("Number of records in 'n' and 't' should be equal: %d != %d", ncount, tcount)
	}
	t.Fatal()

	// err = LoadFile(db, 0, 10, txAbstractDB, os.Getenv("DATADIR") + "/abstracts")
}
