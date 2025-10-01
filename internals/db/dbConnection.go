package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {

	dsn := os.Getenv("DATABASE_URL")
        if dsn == "" {
                dsn = "postgres://datinguser:yourpassword@postgres:5432/datingdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Check if database is reachable
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the database")
	DB = db
	return db, nil
}
