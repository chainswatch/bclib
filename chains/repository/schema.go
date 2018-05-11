package repository

import (
  "github.com/jmoiron/sqlx"
  _ "github.com/lib/pq"
  log "github.com/sirupsen/logrus"
  "fmt"
)

const (
  host     = "watchers_db_1"
  port     = 5432
  user     = "postgres"
  password = "123456"
  dbname   = "postgres"
)

func newPool() *sqlx.DB {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  log.Info(psqlInfo)
  conn, err := sqlx.Open("postgres", psqlInfo)
  if err != nil {
    log.Fatal("m=newPool", err)
  }

  conn.SetMaxIdleConns(0)
  return conn
}

func NewTable(db *sqlx.DB, name string, schema string, drop bool) {
  if drop {
    schema = "DROP TABLE IF EXISTS " + name + " CASCADE;" + schema
  }
  _, err := db.Exec(schema)
  if err != nil {
    log.Warn(name, ": ", err)
  }
}

func CreateBlock(db *sqlx.DB, drop bool) {
  schema := `
  CREATE TABLE blocks (
    created_at timestamp with time zone,
    n_version integer,
    n_height integer NOT NULL UNIQUE,
    n_status bigint,
    n_tx bigint,
    n_file integer,
    n_data_pos bigint,
    n_undo_pos bigint,
    hash_block bytea,
    hash_prev_block bytea,
    hash_merkle_root bytea,
    n_time timestamp with time zone,
    n_bits bigint,
    n_nonce bigint,
    target_difficulty bigint,
    length bigint,
    price numeric(12,4),
    CONSTRAINT blocks_pkey PRIMARY KEY (n_height)
  );`
  NewTable(db, "blocks", schema, drop)
  schema = `
  CREATE TABLE transactions (
    tx_hash bytea NOT NULL UNIQUE,
    n_height integer REFERENCES blocks(n_height) ON DELETE CASCADE,
    n_version integer,
    locktime bigint
  );`
  NewTable(db, "transactions", schema, drop)
  schema = `
  CREATE TABLE tx_inputs (
    tx_hash bytea,
    hash bytea,
    index bigint,
    script bytea,
    sequence bigint
  );`
  NewTable(db, "tx_inputs", schema, drop)
  schema = `
  CREATE TABLE tx_outputs (
    tx_hash bytea,
    value bigint,
    script bytea
  );`
  NewTable(db, "tx_outputs", schema, drop)
}

func Create(drop bool) *sqlx.DB {
  db := newPool()
  err := db.Ping()
  if err != nil {
    log.Fatal(err)
  }
  CreateBlock(db, true)
  schema := `
  CREATE TABLE entities (
    entity_id bigint NOT NULL UNIQUE,
    name char(10)
  );`
  NewTable(db, "entities", schema, drop)
  schema = `
  CREATE TABLE addresses (
    hash bytea,
    entity_id bigint REFERENCES entities(entity_id) ON DELETE CASCADE,
    n_height bigint
  );`
  NewTable(db, "addresses", schema, drop)
  return db
}
