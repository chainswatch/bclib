package repository

import (
  log "github.com/sirupsen/logrus"
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/syndtr/goleveldb/leveldb/opt"
  "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
  "fmt"
)

const (
  host     = "watchers_db_1"
  port     = 5432
  user     = "postgres"
  password = "123456"
  dbname   = "postgres"
)

func ConnectPg() *sqlx.DB {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  log.Info(psqlInfo)
  pg, err := sqlx.Connect("postgres", psqlInfo)
  if err != nil {
    log.Fatal(err)
  }
  return pg
}

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
  if err != nil {
    log.Warn("Error:", err)
    return nil, err
  }
  return db, nil
}
