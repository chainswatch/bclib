package repository

import (
  "database/sql"
  "github.com/jmoiron/sqlx"
  "app/models"
  "fmt"
  log "github.com/sirupsen/logrus"
)
func GetHeaderFromHeight(db *sqlx.DB, nHeight int) models.BlockHeader {
  query := fmt.Sprintf("SELECT n_height, n_file, n_data_pos, n_undo_pos FROM blocks WHERE n_height=$1")
  res := models.BlockHeader{}
  err := db.Get(&res, query, nHeight)
  if err != nil {
    log.Warn(err)
  }
  return res
}

/*
func GetHeaderFromHeight(db *sqlx.DB, nHeight int) (*sql.Row) {
  // query := `SELECT n_height, n_file, n_data_pos, n_undo_pos FROM blocks WHERE n_height = :n_height`
  // return db.NamedQuery(query, m)
  query := fmt.Sprintf(`SELECT n_height, n_file, n_data_pos, n_undo_pos FROM blocks WHERE n_height = ?`, nHeight)
  return db.QueryRow(query)
}
*/

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
	schema := `CREATE TABLE blocks (
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
		CONSTRAINT blocks_pkey PRIMARY KEY (n_height)
	);`
	db.Exec(schema)
}
