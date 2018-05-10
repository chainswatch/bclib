package repository

import (
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/syndtr/goleveldb/leveldb/opt"
)

func GetFlag(db *leveldb.DB, name []byte) (bool, error) {
	command := append([]byte("F"), byte(len(name)))
	command = append(command, name...)
	data, err := db.Get(command, nil)
	if err != nil {
		return false, err
	}

	return data[0] == []byte("1")[0], nil
}

func OpenIndexDb(dataDir string) (*leveldb.DB, error) {
  db, err := leveldb.OpenFile(dataDir + "/blocks/index", &opt.Options{
    ReadOnly: true,
  })
  return db, err
}
