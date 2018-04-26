package db

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "log"
  "fmt"
  "time"
)

type BlockHeader struct {
  gorm.Model
  // ID        uint `gorm:"primary_key"`
  Version        int32
  Height         int32    `gorm:"primary_key;unique"`
  Status         uint32
  NTx            uint32
  NFile          int32
  NDataPos       uint32
  NUndoPos       uint32
  HashPrev       []byte
  HashMerkleRoot []byte
  Timestamp      time.Time
  NBits          uint32
  NNonce         uint32
}

const (
  host     = "watchers_db_1"
  port     = 5432
  user     = "postgres"
  password = "123456"
  dbname   = "postgres"
)

func DbInit() *gorm.DB {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  fmt.Println(psqlInfo)
  db, err := gorm.Open("postgres", psqlInfo)
  if err != nil {
    log.Fatal(err)
  }
  return db
}
