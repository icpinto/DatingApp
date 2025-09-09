package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

func CreateOrUpdateProfile(db *sql.DB, username string, profile models.Profile) error {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		log.Printf("CreateOrUpdateProfile user lookup error for %s: %v", username, err)
		return err
	}
	profile.UserID = userID
	if err := repositories.UpsertProfile(db, profile); err != nil {
		log.Printf("CreateOrUpdateProfile repository error for user %d: %v", userID, err)
		return err
	}
	return nil
}

func GetProfile(db *sql.DB, username string) (models.Profile, error) {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		log.Printf("GetProfile user lookup error for %s: %v", username, err)
		return models.Profile{}, err
	}
	profile, err := repositories.GetProfileByUserID(db, userID)
	if err != nil {
		log.Printf("GetProfile repository error for user %d: %v", userID, err)
		return models.Profile{}, err
	}
	return profile, nil
}

func GetProfiles(db *sql.DB) ([]models.Profile, error) {
	profiles, err := repositories.GetAllProfiles(db)
	if err != nil {
		log.Printf("GetProfiles repository error: %v", err)
		return nil, err
	}
	return profiles, nil
}

func GetProfileByUserID(db *sql.DB, userID int) (models.Profile, error) {
	profile, err := repositories.GetProfileByUserID(db, userID)
	if err != nil {
		log.Printf("GetProfileByUserID repository error for user %d: %v", userID, err)
		return models.Profile{}, err
	}
	return profile, nil
}
