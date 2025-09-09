package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/utils"
)

func GetUsepwd(username string, db *sql.DB) (string, int, error) {
	pwd, id, err := repositories.GetUserpwdByUsername(db, username)
	if err != nil {
		log.Printf("GetUsepwd service error for %s: %v", username, err)
	}
	return pwd, id, err
}

func RegisterUser(db *sql.DB, user models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("RegisterUser hash password error for %s: %v", user.Username, err)
		return err
	}
	user.Password = hashedPassword
	if err := repositories.CreateUser(db, user); err != nil {
		log.Printf("RegisterUser repository error for %s: %v", user.Username, err)
		return err
	}
	return nil
}
