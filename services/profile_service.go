package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

// ProfileService provides operations for user profiles.
type ProfileService struct {
	db *sql.DB
}

// NewProfileService creates a new ProfileService.
func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{db: db}
}

// CreateOrUpdateProfile creates or updates a user's profile.
func (s *ProfileService) CreateOrUpdateProfile(username string, profile models.Profile) error {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("CreateOrUpdateProfile user lookup error for %s: %v", username, err)
		return err
	}
	profile.UserID = userID
	if err := repositories.UpsertProfile(s.db, profile); err != nil {
		log.Printf("CreateOrUpdateProfile repository error for user %d: %v", userID, err)
		return err
	}
	return nil
}

// GetProfile retrieves a user's profile by username.
func (s *ProfileService) GetProfile(username string) (models.Profile, error) {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("GetProfile user lookup error for %s: %v", username, err)
		return models.Profile{}, err
	}
	profile, err := repositories.GetProfileByUserID(s.db, userID)
	if err != nil {
		log.Printf("GetProfile repository error for user %d: %v", userID, err)
		return models.Profile{}, err
	}
	return profile, nil
}

// GetProfiles retrieves all profiles.
func (s *ProfileService) GetProfiles() ([]models.Profile, error) {
	profiles, err := repositories.GetAllProfiles(s.db)
	if err != nil {
		log.Printf("GetProfiles repository error: %v", err)
		return nil, err
	}
	return profiles, nil
}

// GetProfileByUserID retrieves a profile by user ID.
func (s *ProfileService) GetProfileByUserID(userID int) (models.Profile, error) {
	profile, err := repositories.GetProfileByUserID(s.db, userID)
	if err != nil {
		log.Printf("GetProfileByUserID repository error for user %d: %v", userID, err)
		return models.Profile{}, err
	}
	return profile, nil
}
