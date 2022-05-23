package postgres

import (
	"database/sql"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err

	}
	return &DB{
		conn: db,
	}, nil
}

func (d *DB) Store(*domain.URL) error {
	return nil
}

func (d *DB) FindByKey(string) (*domain.URL, error) {
	return nil, nil
}

func (d *DB) FindAll(string) []*domain.URL {
	return nil
}

func (d *DB) Close() error {
	err := d.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) Ping() error {
	err := d.conn.Ping()
	if err != nil {
		return err
	}
	return nil
}
