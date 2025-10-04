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
	err := db.QueryRow("SELECT id FROM users WHERE username=$1 AND is_active = true", username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		log.Printf("GetUserIDByUsername query error for %s: %v", username, err)
		return 0, err
	}
	return id, nil
}

func GetUsernameByID(db *sql.DB, userID int) (string, error) {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id=$1 AND is_active = true", userID).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}

		log.Printf("GetUsernameByID query error for %d: %v", userID, err)

		return "", err
	}
	return username, nil
}

// GetUserStatusByID returns whether the specified user is currently active.
func GetUserStatusByID(db *sql.DB, userID int) (bool, error) {
	var isActive bool
	err := db.QueryRow("SELECT is_active FROM users WHERE id=$1", userID).Scan(&isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrUserNotFound
		}
		log.Printf("GetUserStatusByID query error for %d: %v", userID, err)
		return false, err
	}
	return isActive, nil
}

// DeactivateUserTx sets a user's account as inactive within the supplied transaction.
func DeactivateUserTx(tx *sql.Tx, userID int) error {
	res, err := tx.Exec(`
        UPDATE users
        SET is_active = false, deactivated_at = NOW()
        WHERE id = $1 AND is_active = true`, userID)
	if err != nil {
		log.Printf("DeactivateUserTx exec error for user %d: %v", userID, err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

// ReactivateUserTx sets a user's account as active within the supplied transaction.
func ReactivateUserTx(tx *sql.Tx, userID int) error {
	res, err := tx.Exec(`
        UPDATE users
        SET is_active = true, deactivated_at = NULL
        WHERE id = $1 AND is_active = false`, userID)
	if err != nil {
		log.Printf("ReactivateUserTx exec error for user %d: %v", userID, err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

// DeleteUserTx hard deletes a user row from the database within the supplied transaction.
func DeleteUserTx(tx *sql.Tx, userID int) error {
	res, err := tx.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		log.Printf("DeleteUserTx exec error for user %d: %v", userID, err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}
