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
	fileNum				uint32
	NBlocks				uint32
	NSize					uint32	
	NUndoSize			uint32	
	NHeightFirst	uint32	
	NHeightLast		uint32	
	NTimeFirst		uint32	
	NTimeLast			uint32	
}

func decodeBlockFileIdx(br parser.Reader) *blockFile {
	f := &blockFile{}
	f.NBlocks = uint32(br.ReadVarint())
	f.NSize = uint32(br.ReadVarint())
	f.NUndoSize = uint32(br.ReadVarint())
	f.NHeightFirst = uint32(br.ReadVarint())
	f.NHeightLast = uint32(br.ReadVarint())
	f.NTimeFirst =  uint32(br.ReadVarint())
	f.NTimeLast =  uint32(br.ReadVarint())
	return f
}

// Get block index record by hash
func blockIndexRecord(db *leveldb.DB, h []byte) (bh models.BlockHeader, err error) {
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

// OpenIndexDb gets transaction index record
func OpenIndexDb() (*leveldb.DB, error) {
	dataDir := os.Getenv("DATADIR")
  db, err := leveldb.OpenFile(dataDir + "/blocks/index", &opt.Options{
    ReadOnly: true,
  })
  return db, err
}

/*
func GetBlockHeaders() {
  var err error
  indexDb, err = db.OpenIndexDb(btc.DataDir) // TODO: Error handling
  if err != nil {
    log.Warn("Error:", err)
  }
  defer indexDb.Close()

  iter := indexDb.NewIterator(util.BytesPrefix([]byte("b")), nil)
  for iter.Next() {
    btc.HashBlock = iter.Key()
    data := iter.Value()
    // log.Info(data)
    btc.parseBlockHeaderData(data)
    // Copy, then insert in DB
    // _, err = db.InsertHeader(btc.SqlDb, btc.BlockHeader)
    if err != nil {
      if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
        log.Warn(err)
      }
    }
  }
  iter.Release()
  err = iter.Error()
  if err != nil {
    log.Warn(err)
  }
}
*/


/*
func getTransaction() {
  f, err := db.GetFlag(btc.IndexDb, []byte("txindex"))
  if err != nil {
    log.Warn(err)
  }
  if !f {
    fmt.Println("txindex is not enabled for your bitcoind")
  }

  result, err := db.GetTxIndexRecordByBigEndianHex(indexDb, args[1])
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", result)

  tx, err := blockchainparser.NewTxFromFile(datadir, magicId, uint32(result.NFile), result.NDataPos, result.NTxOffset)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%+v\n", tx)
}
*/
