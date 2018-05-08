package repository

import (
  "github.com/jmoiron/sqlx"
  log "github.com/sirupsen/logrus"
)

func NewTable(db *sqlx.DB, schema string) {
  _, err := db.Exec(schema)
  if err != nil {
    log.Warn(err)
  }
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
  NewTable(db, schema)
	schema = `
  DROP TABLE IF EXISTS transactions CASCADE;
  CREATE TABLE transactions (
    tx_hash bytea NOT NULL UNIQUE,
    n_height integer REFERENCES blocks(n_height) ON DELETE CASCADE,
    n_version integer,
    locktime bigint
	);`
  NewTable(db, schema)
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
  NewTable(db, schema)
	schema = `
  DROP TABLE IF EXISTS tx_outputs;
  CREATE TABLE tx_outputs (
    tx_hash bytea REFERENCES transactions(tx_hash) ON DELETE CASCADE,
    value bigint,
    script bytea,
    PRIMARY KEY (tx_hash)
	);`
  NewTable(db, schema)
	schema = `
  DROP TABLE IF EXISTS entities;
  CREATE TABLE entities (
    entity_id bigint NOT NULL UNIQUE,
    name char(10)
	);`
  NewTable(db, schema)
	schema = `
  DROP TABLE IF EXISTS addresses;
  CREATE TABLE addresses (
    hash bytea,
    entity_id bigint REFERENCES entities(entity_id) ON DELETE CASCADE,
    n_height bigint
	);`
  NewTable(db, schema)
}
