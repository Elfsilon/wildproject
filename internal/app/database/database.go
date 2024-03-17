package db

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var (
	ErrDatabaseAlreadyOpened = errors.New("database already opened")
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Open(conn string) error {
	if d.DB != nil {
		return ErrDatabaseAlreadyOpened
	}

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	d.DB = db
	return nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
