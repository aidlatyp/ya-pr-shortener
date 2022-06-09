package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dsn string) (*DB, error) {
	if dsn == "" {
		return nil, errors.New("invalid connection string")
	}
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db := DB{conn: conn}
	db.createTablesIfNotExits()

	return &db, nil
}

func (d *DB) BatchDelete(urls []string, id string) error {

	query := `UPDATE public.urls SET del = true WHERE id = any ($1) AND user_id=$2;`

	stmt, err := d.conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("error while prepare stmt: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}(stmt)

	if _, err = stmt.Exec(pq.Array(urls), id); err != nil {
		return err
	}
	return nil
}

func (d *DB) BatchWrite(uris []domain.URL) error {

	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}

	query := `SELECT id FROM public.users WHERE  id = $1 `
	row := d.conn.QueryRow(query, uris[0].Owner)

	var s string
	haveRows := row.Scan(&s)
	if haveRows == sql.ErrNoRows {
		insert := `INSERT INTO public.users (id) VALUES ($1);`
		_, err := tx.Exec(insert, uris[0].Owner)
		if err != nil {
			return errors.New("error while trying insert user")
		}
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO urls (id, orig_url, user_id) VALUES ($1,$2,$3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, u := range uris {
		if _, err = stmt.Exec(u.Short, u.Orig, u.Owner); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DB) Store(url *domain.URL) error {

	// todo wrap into TX
	query := `SELECT id FROM public.users WHERE  id = $1`
	row := d.conn.QueryRow(query, url.Owner)
	var s string
	haveRows := row.Scan(&s)
	if haveRows == sql.ErrNoRows {
		insert := `INSERT INTO public.users (id) VALUES ($1);`
		_, err := d.conn.Exec(insert, url.Owner)
		if err != nil {
			log.Println(err.Error())
			return errors.New("error while trying insert user")
		}
	}

	prep, err := d.conn.Prepare("INSERT INTO public.urls (id, orig_url, user_id)" +
		" VALUES ($1, $2, $3) ON CONFLICT (user_id,orig_url) DO UPDATE SET orig_url=EXCLUDED.orig_url  RETURNING id")
	if err != nil {
		return err
	}

	result := prep.QueryRowContext(context.Background(), url.Short, url.Orig, url.Owner)

	var id string
	result.Scan(&id)

	if id != url.Short {
		return usecase.ErrAlreadyExists{
			Err:            errors.New("duplicate entry, given entity record already exists"),
			ExistShortenID: id,
			Orig:           url.Orig,
		}
	}
	return nil
}

func (d *DB) FindByKey(key string) (*domain.URL, error) {

	var isDel bool
	query := `SELECT id,orig_url,user_id, del FROM public.urls WHERE id = $1;`

	row := d.conn.QueryRow(query, key)
	url := domain.URL{}
	err := row.Scan(&url.Short, &url.Orig, &url.Owner, &isDel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("key %v not exists", key)
		}
		return nil, err
	}

	if isDel {
		return nil, usecase.ErrURLDeleted{
			Err:     errors.New("requested url wad deleted"),
			ShortID: key,
		}
	}
	return &url, nil
}

func (d *DB) FindAll(key string) []*domain.URL {

	result := make([]*domain.URL, 0)
	query := "SELECT id,orig_url, user_id FROM public.urls WHERE user_id = $1;"

	rows, err := d.conn.Query(query, key)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		url := domain.URL{}
		err = rows.Scan(&url.Short, &url.Orig, &url.Owner)
		if err != nil {
			if err == sql.ErrNoRows {
				return result
			}
			return []*domain.URL{}
		}
		result = append(result, &url)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return result
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) Ping() error {
	return d.conn.Ping()
}

func (d *DB) createTablesIfNotExits() {
	_, err := d.conn.Exec(`CREATE TABLE IF NOT EXISTS public.users (
    											id TEXT NOT NULL,
												CONSTRAINT user_constraint PRIMARY KEY (id));

						   CREATE TABLE IF NOT EXISTS public.urls (
                                 id TEXT NOT NULL,
                                 orig_url TEXT NOT NULL,
                                 user_id TEXT NOT NULL,
                                 del  BOOLEAN DEFAULT false NOT NULL,
                                 CONSTRAINT url_constraint PRIMARY KEY (id),
                                 CONSTRAINT orig_url_constraint UNIQUE (user_id, orig_url),
                                 FOREIGN KEY (user_id) REFERENCES public.users (id));`)
	if err != nil {
		log.Println(err.Error())
	}
}
