package repositories

import (
	"database/sql"
)

func GetUserpwdByUsername(db *sql.DB, username string) (string, int, error) {
	var hashedPassword string
	var userId int

	if err := db.QueryRow("SELECT id, password FROM users WHERE username=$1", username).
		Scan(&userId, &hashedPassword); err != nil {
		return "", 0, err
	}
	return hashedPassword, userId, nil

}
