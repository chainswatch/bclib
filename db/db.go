package db

import (
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/syndtr/goleveldb/leveldb/opt"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres" // trivial postgres interactions
  "log"
  "fmt"
)

const (
  host     = "watchers_db_1"
  port     = 5432
  user     = "postgres"
  password = "123456"
  dbname   = "postgres"
)

// Init connects to Postgres DB
func PgInit() *gorm.DB {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  fmt.Println(psqlInfo)
  db, err := gorm.Open("postgres", psqlInfo)
  if err != nil {
    log.Fatal(err)
  }
  return db
}

func OpenIndexDb(dataDir string) (*leveldb.DB, error) {
  db, err := leveldb.OpenFile(dataDir + "/blocks/index", &opt.Options{
    ReadOnly: true,
  })
  if err != nil {
    fmt.Println("Error:", err)
    return nil, err
  }
  return db, nil
}
