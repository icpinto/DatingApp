package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

// ProfileService provides operations for user profiles.
type ProfileService struct {
	db   *sql.DB
	repo *repositories.ProfileRepository
}

// NewProfileService creates a new ProfileService.
func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{db: db, repo: repositories.NewProfileRepository(db)}
}

// CreateOrUpdateProfile creates or updates a user's profile.
func (s *ProfileService) CreateOrUpdateProfile(username string, profile models.Profile) error {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("CreateOrUpdateProfile user lookup error for %s: %v", username, err)
		return err
	}
	profile.UserID = userID
	if err := s.repo.Upsert(profile); err != nil {
		log.Printf("CreateOrUpdateProfile repository error for user %d: %v", userID, err)
		return err
	}
	return nil
}

// GetProfile retrieves a user's profile by username.
func (s *ProfileService) GetProfile(username string) (models.UserProfile, error) {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("GetProfile user lookup error for %s: %v", username, err)
		return models.UserProfile{}, err
	}
	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		log.Printf("GetProfile repository error for user %d: %v", userID, err)
		return models.UserProfile{}, err
	}
	return profile, nil
}

// GetProfiles retrieves all profiles.
func (s *ProfileService) GetProfiles() ([]models.UserProfile, error) {
	profiles, err := s.repo.GetAll()
	if err != nil {
		log.Printf("GetProfiles repository error: %v", err)
		return nil, err
	}
	return profiles, nil
}

// GetProfileByUserID retrieves a profile by user ID.
func (s *ProfileService) GetProfileByUserID(userID int) (models.UserProfile, error) {
	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		log.Printf("GetProfileByUserID repository error for user %d: %v", userID, err)
		return models.UserProfile{}, err
	}
	return profile, nil
}

// GetProfileEnums returns available enum options for profiles.
func (s *ProfileService) GetProfileEnums() (models.ProfileEnums, error) {
	enums, err := s.repo.GetProfileEnums()
	if err != nil {
		log.Printf("GetProfileEnums repository error: %v", err)
		return models.ProfileEnums{}, err
	}
	return enums, nil
}
