package database

import (
	"database/sql"
	"errors"
)

var (
	ErrAlreadyOpened = errors.New("cannot open: database connection already opened")
	ErrNotOpened     = errors.New("cannot close: no opened database connections found")
)

type Database interface {
	Open(conn string) error
	Close() error
}

type Instance interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
