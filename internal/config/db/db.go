package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func CreateTableIfNotExist(db *sql.DB) error {
	
	query := `CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price REAL NOT NULL,
	sizes TEXT NOT NULL,
    images TEXT NOT NULL,
    is_avaliable BOOLEAN DEFAULT TRUE
);
`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
	return nil
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}