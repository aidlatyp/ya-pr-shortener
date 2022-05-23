package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dsn string) (*DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db := DB{conn: conn}
	db.init()
	return &db, nil
}

func (d *DB) Store(url *domain.URL) error {
	// todo TX
	q := fmt.Sprintf(`SELECT id FROM public.users WHERE  id = '%s';`, url.Owner)
	row := d.conn.QueryRow(q)
	var s string
	haveRows := row.Scan(&s)
	if haveRows == sql.ErrNoRows {
		insert := fmt.Sprintf(`INSERT INTO public.users (id) VALUES ('%s');`, url.Owner)
		_, err := d.conn.Exec(insert)
		if err != nil {
			log.Println(err.Error())
			return errors.New("error while trying insert user")
		}
	}

	insert := fmt.Sprintf(`INSERT INTO "public"."urls"(id,        orig_url, user_id)
						VALUES ('%s','%s','%s')`, url.Short, url.Orig, url.Owner)
	_, err := d.conn.Exec(insert)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error while trying insert url")
	}
	return nil
}

func (d *DB) FindByKey(id string) (*domain.URL, error) {

	q := fmt.Sprintf(`SELECT id,orig_url, user_id FROM public.urls WHERE id = '%s';`, id)
	row := d.conn.QueryRow(q)
	url := domain.URL{}
	err := row.Scan(&url.Short, &url.Orig, &url.Owner)
	if err != nil {
		fmt.Println("NIL")
		if err == sql.ErrNoRows {
			fmt.Println("NOT FOUND")
			return nil, fmt.Errorf("key %v not exists", id)
		}
		return nil, err
	}
	return &url, nil
}

func (d *DB) FindAll(user string) []*domain.URL {

	result := make([]*domain.URL, 0)

	q := fmt.Sprintf(`SELECT id,orig_url, user_id FROM public.urls WHERE user_id = '%s';`, user)
	rows, err := d.conn.Query(q)
	defer rows.Close()
	if err != nil {
		return result
	}

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
	return result
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

func (d *DB) init() {
	_, err := d.conn.Exec(`CREATE TABLE IF NOT EXISTS public.users (
    											id TEXT NOT NULL,
												CONSTRAINT user_constraint PRIMARY KEY (id));

						   CREATE TABLE IF NOT EXISTS public.urls (
                                 id TEXT NOT NULL,
                                 orig_url TEXT NOT NULL,
                                 user_id TEXT NOT NULL,
                                 CONSTRAINT url_constraint PRIMARY KEY (id),
                                 FOREIGN KEY (user_id) REFERENCES public.users (id));`)

	if err != nil {
		log.Println(err.Error())
	}
}
