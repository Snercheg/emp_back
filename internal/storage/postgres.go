package storage

import (
	"database/sql"
	"fmt"
)

type DBConfig struct {
	DBUrl string
}
type Storage struct {
	DB     *sql.DB
	config DBConfig
}

func New(cfg DBConfig) (*sql.DB, error) {
	op := "storage.postgres.New"
	db, err := sql.Open("pgx", cfg.DBUrl)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	return db, nil
}
