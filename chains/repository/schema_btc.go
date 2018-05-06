package repository

import (
  "database/sql"
  "github.com/jmoiron/sqlx"
  "app/models"
  "fmt"
  log "github.com/sirupsen/logrus"
)

func GetRowCount(db *sqlx.DB, table string) (int, error) {
  var id int
  query := fmt.Sprintf("SELECT count(*) FROM %s", table)
  err := db.Get(&id, query)
  return id, err
}

func GetHeaderFromHeight(db *sqlx.DB, nHeight int) (models.BlockHeader, error) {
  query := fmt.Sprintf("SELECT n_height, n_file, n_data_pos, n_undo_pos, hash_block FROM blocks WHERE n_height=$1")
  res := models.BlockHeader{}
  err := db.Get(&res, query, nHeight)
  return res, err
}

func InsertTransaction(db *sqlx.DB, m models.Transaction, hash_block models.Hash256) (*sql.Rows, error) {
  query := `
  INSERT INTO transactions (
    tx_hash,
    hash_block,
    n_version,
    locktime
  ) VALUES (
    $1,
    $2,
    $3,
    $4
  )
  `
  return db.Query(query, m.Hash, hash_block, m.NVersion, m.Locktime)
}

func InsertInput(db *sqlx.DB, m models.TxInput, tx_hash models.Hash256) (*sql.Rows, error) {
  query := `
  INSERT INTO tx_inputs (
    tx_hash,
    hash,
    index,
    script,
    sequence
  ) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
  )
  `
  return db.Query(query, tx_hash, m.Hash, m.Index, m.Script, m.Sequence)
}

func InsertOutput(db *sqlx.DB, m models.TxOutput, tx_hash models.Hash256) (*sql.Rows, error) {
  query := `
  INSERT INTO tx_outputs (
    tx_hash,
    value,
    script
  ) VALUES (
    $1,
    $2,
    $3
  )
  `
  return db.Query(query, tx_hash, m.Value, m.Script)
}


func InsertHeader(db *sqlx.DB, m models.BlockHeader) (sql.Result, error) {
  query := `INSERT INTO blocks (
		n_version,
		n_height,
		n_status,
		n_tx,
		n_file,
		n_data_pos,
		n_undo_pos,
		hash_block,
		hash_prev_block,
		hash_merkle_root,
		n_time,
		n_bits,
		n_nonce
		) VALUES (
      :n_version,
      :n_height,
      :n_status,
      :n_tx,
      :n_file,
      :n_data_pos,
      :n_undo_pos,
      :hash_block,
      :hash_prev_block,
      :hash_merkle_root,
      :n_time,
      :n_bits,
      :n_nonce
    )`
    return db.NamedExec(query, m)
}

func CreateBtc() {
	db := ConnectPg()
  defer db.Close()
	schema := `
  DROP TABLE IF EXISTS blocks CASCADE;
  CREATE TABLE blocks (
		created_at timestamp with time zone,
		n_version integer,
		n_height integer NOT NULL UNIQUE,
		n_status bigint,
		n_tx bigint,
		n_file integer,
		n_data_pos bigint,
		n_undo_pos bigint,
		hash_block bytea UNIQUE,
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
  _, err := db.Exec(schema)
  if err != nil {
    log.Warn(err)
  }
	schema = `
  DROP TABLE IF EXISTS transactions CASCADE;
  CREATE TABLE transactions (
    tx_hash bytea NOT NULL UNIQUE,
    hash_block bytea REFERENCES blocks(hash_block) ON DELETE CASCADE,
    n_version integer,
    locktime bigint
	);`
  _, err = db.Exec(schema)
  if err != nil {
    log.Warn(err)
  }
	schema = `
  DROP TABLE IF EXISTS tx_inputs;
  CREATE TABLE tx_inputs (
    tx_hash bytea REFERENCES transactions(tx_hash) ON DELETE CASCADE,
    hash bytea,
    index bigint,
    script bytea,
    sequence bigint,
    PRIMARY KEY (tx_hash)
	);`
	db.Exec(schema)
  _, err = db.Exec(schema)
  if err != nil {
    log.Warn(err)
  }
	schema = `
  DROP TABLE IF EXISTS tx_outputs;
  CREATE TABLE tx_outputs (
    tx_hash bytea REFERENCES transactions(tx_hash) ON DELETE CASCADE,
    value bigint,
    script bytea,
    PRIMARY KEY (tx_hash)
	);`
	db.Exec(schema)
  _, err = db.Exec(schema)
  if err != nil {
    log.Warn(err)
  }
}
