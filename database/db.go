package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func OpenDB() *Database {
	db, err := sql.Open("mysql", "root:YsRZeAQVoFtXgaYfkkNEtuxQiPbcqkMH@yamabiko.proxy.rlwy.net:42537/railway")
	if err != nil {
		panic(err)
	}
	return &Database{
		db: db,
	}
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func (d *Database) CloseDB() {
	d.db.Close()
}
