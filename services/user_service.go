package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/utils"
)

// UserService provides user related operations.
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new UserService with the given database handle.
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// GetUsepwd retrieves the hashed password and user ID for a username.
func (s *UserService) GetUsepwd(username string) (string, int, error) {
	pwd, id, err := repositories.GetUserpwdByUsername(s.db, username)
	if err != nil {
		log.Printf("GetUsepwd service error for %s: %v", username, err)
	}
	return pwd, id, err
}

// RegisterUser creates a new user in the database after hashing the password.
func (s *UserService) RegisterUser(user models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("RegisterUser hash password error for %s: %v", user.Username, err)
		return err
	}
	user.Password = hashedPassword
	if err := repositories.CreateUser(s.db, user); err != nil {
		log.Printf("RegisterUser repository error for %s: %v", user.Username, err)
		return err
	}
	return nil
}

// GetUserIDByUsername returns the user ID for a given username.
func (s *UserService) GetUserIDByUsername(username string) (int, error) {
	id, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("GetUserIDByUsername service error for %s: %v", username, err)
	}
	return id, err
}
