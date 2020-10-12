package db

import (
	"database/sql"
)

type DB struct {
	*sql.DB
}

const (
	openConns = 25
	idleConns = 25
)

func NewDB(source string) (*DB, error) {
	db, err := sql.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(openConns)
	db.SetMaxIdleConns(idleConns)

	return &DB{db}, nil
}

func (db *DB) IsNoRowsErr(err error) bool {
	return err == sql.ErrNoRows
}
