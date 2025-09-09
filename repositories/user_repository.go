package repositories

import (
	"database/sql"

	"github.com/icpinto/dating-app/models"
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

func CreateUser(db *sql.DB, user models.User) error {
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	return err
}

func GetUserIDByUsername(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&id)
	return id, err
}
