package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	instance *sql.DB
}

func NewPostgres() Database {
	return &Postgres{}
}

func (pg *Postgres) Instance() (Instance, error) {
	if pg.instance == nil {
		return nil, ErrNotOpened
	}

	return pg.instance, nil
}

func (pg *Postgres) Open(conn string) error {
	if pg.instance != nil {
		return ErrAlreadyOpened
	}

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	pg.instance = db
	return nil
}

func (pg *Postgres) Close() error {
	if pg.instance == nil {
		return ErrNotOpened
	}

	return pg.instance.Close()
}
