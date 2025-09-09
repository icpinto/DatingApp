package services

import (
	"database/sql"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

func CreateOrUpdateProfile(db *sql.DB, username string, profile models.Profile) error {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		return err
	}
	profile.UserID = userID
	return repositories.UpsertProfile(db, profile)
}

func GetProfile(db *sql.DB, username string) (models.Profile, error) {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		return models.Profile{}, err
	}
	return repositories.GetProfileByUserID(db, userID)
}

func GetProfiles(db *sql.DB) ([]models.Profile, error) {
	return repositories.GetAllProfiles(db)
}

func GetProfileByUserID(db *sql.DB, userID int) (models.Profile, error) {
	return repositories.GetProfileByUserID(db, userID)
}
