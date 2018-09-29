package repository

import (
  log "github.com/sirupsen/logrus"
  "database/sql"
  "github.com/jmoiron/sqlx"
  "app/models"
  "fmt"
)

func GetLastBlockHeader(db *sqlx.DB) (models.BlockHeader, error) {
  query := fmt.Sprintf("SELECT n_height, n_file, n_data_pos, length FROM blocks ORDER BY n_height DESC LIMIT 1")
  // query := fmt.Sprintf("SELECT * FROM (SELECT n_height, n_file, n_data_pos, length FROM blocks ORDER BY n_height DESC LIMIT 2) table_alias ORDER BY n_height LIMIT 1")
  row := db.QueryRow(query)
  res := models.BlockHeader{}
  err := row.Scan(&res.NHeight, &res.NFile, &res.NDataPos, &res.Length)
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

func PrepareInsertTransaction(tx *sql.Tx) func(m models.Transaction) {
  query := `
  INSERT INTO transactions (
    tx_hash, n_version, n_vin, n_vout, locktime
  ) VALUES (
    $1, $2, $3, $4, $5
  )
  `
  stmt, err := tx.Prepare(query)
  if err != nil {
    log.Panic(err)
  }

  return func(m models.Transaction) {
    _, err := stmt.Exec(m.Hash, m.NVersion, m.NVin, m.NVout, m.Locktime)
    if err != nil {
      log.Warn(err)
      log.Info(m.Hash)
    }
  }
}

func PrepareInsertInput(tx *sql.Tx) func(m models.TxInput, tx_hash models.Hash256) {
  query := `
  INSERT INTO tx_inputs (
    tx_hash, hash, index, sequence
  ) VALUES (
    $1, $2, $3, $4
  )
  `
  stmt, err := tx.Prepare(query)
  if err != nil {
    log.Panic("PrepareInsertInput: ", err)
  }

  return func(m models.TxInput, tx_hash models.Hash256) {
    _, err := stmt.Exec(tx_hash, m.Hash, m.Index, m.Sequence)
    if err != nil {
      log.Warn("InsertInput: ", err)
    }
  }
}

func PrepareInsertOutput(tx *sql.Tx) func(m models.TxOutput, tx_hash models.Hash256) {
  query := `
  INSERT INTO tx_outputs (
    tx_hash, index, value, hash160, script
  ) VALUES (
    $1, $2, $3, $4, $5
  )
  `
  stmt, err := tx.Prepare(query)
  if err != nil {
    log.Panic("PrepareInsertOutput :", err)
  }

  return func(m models.TxOutput, tx_hash models.Hash256) {
    _, err := stmt.Exec(tx_hash, m.Index, m.Value, m.Hash160, m.Script)
    if err != nil {
      log.Warn("InsertOutput: ", err)
    }
  }
}

func InsertHeader(tx *sql.Tx, m models.Block) {
  query := `INSERT INTO blocks (
    n_version, n_height, n_status, n_tx, n_file,
    n_data_pos, n_undo_pos, hash_block, hash_prev_block, hash_merkle_root,
    n_time, n_bits, n_nonce, target_difficulty, length
  ) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15
  )`

  _, err := tx.Exec(query,
  m.NVersion, m.NHeight, m.NStatus, m.NTx, m.NFile,
  m.NDataPos, m.NUndoPos, m.HashBlock, m.HashPrevBlock, m.HashMerkleRoot,
  m.NTime, m.NBits, m.NNonce, m.TargetDifficulty, m.Length)
  if err != nil {
    log.Warn(err)
  }
}
