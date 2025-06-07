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

CREATE TABLE IF NOT EXISTS address (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    neighborhood TEXT NOT NULL,
    street TEXT NOT NULL,
    number TEXT NOT NULL,
    state TEXT NOT NULL,
    country TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS payment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    email TEXT NOT NULL,
    telephone TEXT NOT NULL,
    cpf TEXT NOT NULL,
    address_id INTEGER NOT NULL,
    FOREIGN KEY (address_id) REFERENCES address(id)
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
