package repository

import (
  log "github.com/sirupsen/logrus"
  "database/sql"
  "github.com/jmoiron/sqlx"
  "app/models"
  "fmt"
)

func GetLastBlockHeader(db *sqlx.DB) (models.BlockHeader, error) {
  query := fmt.Sprintf("SELECT n_height, n_file, n_data_pos FROM blocks ORDER BY n_height DESC LIMIT 1")
  row := db.QueryRow(query)
  res := models.BlockHeader{}
  err := row.Scan(&res)
  return res, err
}

func GetRowCount(db *sqlx.DB, table string) (int, error) {
  var id int
  query := fmt.Sprintf("SELECT count(*) FROM %s", table)
  err := db.Get(&id, query)
  return id, err
}

func GetHeaderFromHeight(db *sqlx.DB, nHeight int) (models.BlockHeader, error) {
  query := fmt.Sprintf("SELECT n_file, n_data_pos FROM blocks WHERE n_height=$1")
  res := models.BlockHeader{}
  err := db.Get(&res, query, nHeight)
  return res, err
}

func InsertTransaction(tx *sql.Tx, m models.Transaction, n_height uint32) {
  query := `
  INSERT INTO transactions (
    tx_hash,
    n_height,
    n_version,
    locktime
  ) VALUES (
    $1,
    $2,
    $3,
    $4
  )
  `
  _, err := tx.Exec(query, m.Hash, n_height, m.NVersion, m.Locktime)
  if err != nil {
    log.Fatal(err)
  }
}

func InsertInput(tx *sql.Tx, m models.TxInput, tx_hash models.Hash256) {
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
  _, err := tx.Exec(query, tx_hash, m.Hash, m.Index, m.Script, m.Sequence)
  if err != nil {
    log.Fatal(err)
  }
}

func InsertOutput(tx *sql.Tx, m models.TxOutput, tx_hash models.Hash256) {
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
  _, err := tx.Exec(query, tx_hash, m.Value, m.Script)
  if err != nil {
    log.Fatal(err)
  }
}

func InsertHeader(tx *sql.Tx, m models.BlockHeader) {
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
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13
  )`
  _, err := tx.Exec(query, m.NVersion, m.NHeight, m.NStatus, m.NTx,
    m.NFile, m.NDataPos, m.NUndoPos, m.HashBlock, m.HashPrevBlock,
    m.HashMerkleRoot, m.NTime, m.NBits, m.NNonce)
  if err != nil {
    log.Fatal(err)
  }
}
