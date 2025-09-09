package services

import (
	"database/sql"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/utils"
)

func GetUsepwd(username string, db *sql.DB) (string, int, error) {
	return repositories.GetUserpwdByUsername(db, username)
}

func RegisterUser(db *sql.DB, user models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return repositories.CreateUser(db, user)
}
