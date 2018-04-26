package db

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "log"
)

func DbInit() {
  db, err := sql.Open("sqlite3", "./data/foo.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    log.Fatal(err)
  }
}
