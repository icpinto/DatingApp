package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {

	db, err := sql.Open("postgres", "user=datinguser dbname=datingapp sslmode=disable password=yourpassword")
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
