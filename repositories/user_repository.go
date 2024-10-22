package repositories

import (
	"database/sql"
)

func GetUserpwdByUsername(db *sql.DB, username string) (string, error) {
	var hashedPassword string

	if err := db.QueryRow("SELECT password FROM users WHERE username=$1", username).
		Scan(&hashedPassword); err != nil {
		return "", err
	}
	return hashedPassword, nil

}
