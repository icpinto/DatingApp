package repositories

import (
	"database/sql"
	"errors"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/lib/pq"
)

var (
	ErrDuplicateUser = errors.New("duplicate user")
	ErrUserNotFound  = errors.New("user not found")
)

func GetUserpwdByUsername(db *sql.DB, username string) (string, int, error) {
	var hashedPassword string
	var userId int

	err := db.QueryRow("SELECT id, password FROM users WHERE username=$1", username).
		Scan(&userId, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", 0, ErrUserNotFound
		}
		log.Printf("GetUserpwdByUsername query error for %s: %v", username, err)
		return "", 0, err
	}
	return hashedPassword, userId, nil
}

func CreateUser(db *sql.DB, user models.User) error {
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			log.Printf("CreateUser duplicate user %s: %v", user.Username, err)
			return ErrDuplicateUser
		}
		log.Printf("CreateUser exec error for %s: %v", user.Username, err)
		return err
	}
	return nil
}

func GetUserIDByUsername(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		log.Printf("GetUserIDByUsername query error for %s: %v", username, err)
		return 0, err
	}
	return id, nil
}

func GetUsernameByID(db *sql.DB, id int) (string, error) {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id=$1", id).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		log.Printf("GetUsernameByID query error for %d: %v", id, err)
		return "", err
	}
	return username, nil
}
