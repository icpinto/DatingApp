package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
	"github.com/lib/pq"
)

// ProfileService provides operations for user profiles.
type ProfileService struct {
	db                *sql.DB
	repo              *repositories.ProfileRepository
	profileOutboxRepo *repositories.ProfileSyncOutboxRepository
}

// NewProfileService creates a new ProfileService.
func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{
		db:                db,
		repo:              repositories.NewProfileRepository(db),
		profileOutboxRepo: repositories.NewProfileSyncOutboxRepository(db),
	}
}

// ErrInvalidEnum indicates an invalid enum value was provided.
var ErrInvalidEnum = errors.New("invalid enum value")

var ErrInvalidVerificationToken = errors.New("invalid verification token")

var ErrVerificationMismatch = errors.New("verification data mismatch")

// CreateOrUpdateProfile creates or updates a user's profile.
func (s *ProfileService) CreateOrUpdateProfile(username string, profile models.Profile, phoneNumber, contactToken, identityToken string) (models.Profile, error) {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("CreateOrUpdateProfile user lookup error for %s: %v", username, err)
		return models.Profile{}, err
	}
	profile.UserID = userID

	existingStatus, err := s.repo.GetVerificationStatus(userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("CreateOrUpdateProfile verification lookup error for user %d: %v", userID, err)
		return models.Profile{}, err
	}

	if phoneNumber != "" {
		profile.PhoneNumber = phoneNumber
	} else if existingStatus.PhoneNumber != "" {
		profile.PhoneNumber = existingStatus.PhoneNumber
	}

	profile.ContactVerified = existingStatus.ContactVerified
	profile.IdentityVerified = existingStatus.IdentityVerified

	if contactToken != "" {
		claims, err := parseVerificationToken(contactToken, getVerificationSecret("CONTACT_VERIFICATION_JWT_SECRET"))
		if err != nil {
			log.Printf("CreateOrUpdateProfile contact token error for user %d: %v", userID, err)
			return models.Profile{}, err
		}
		phoneFromToken := extractPhoneNumber(claims)
		if phoneFromToken != "" {
			if profile.PhoneNumber != "" && !strings.EqualFold(profile.PhoneNumber, phoneFromToken) {
				log.Printf("CreateOrUpdateProfile phone mismatch for user %d", userID)
				return models.Profile{}, ErrVerificationMismatch
			}
			profile.PhoneNumber = phoneFromToken
		} else if profile.PhoneNumber == "" {
			log.Printf("CreateOrUpdateProfile contact token missing phone for user %d", userID)
			return models.Profile{}, ErrInvalidVerificationToken
		}
		profile.ContactVerified = true
	} else if phoneNumber != "" && !strings.EqualFold(phoneNumber, existingStatus.PhoneNumber) {
		profile.ContactVerified = false
	} else if phoneNumber == "" && existingStatus.PhoneNumber != "" {
		profile.ContactVerified = false
	}

	if identityToken != "" {
		if _, err := parseVerificationToken(identityToken, getVerificationSecret("IDENTITY_VERIFICATION_JWT_SECRET")); err != nil {
			log.Printf("CreateOrUpdateProfile identity token error for user %d: %v", userID, err)
			return models.Profile{}, err
		}
		profile.IdentityVerified = true
	}

	profile.Verified = profile.ContactVerified && profile.IdentityVerified

	if err := s.repo.Upsert(profile); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "22P02" && strings.Contains(pqErr.Message, "invalid input value for enum") {
				log.Printf("CreateOrUpdateProfile invalid enum for user %d: %v", userID, pqErr)
				return models.Profile{}, ErrInvalidEnum
			}
		}
		log.Printf("CreateOrUpdateProfile repository error for user %d: %v", userID, err)
		return models.Profile{}, err
	}
	saved, err := s.repo.GetByUserID(userID)
	if err != nil {
		log.Printf("CreateOrUpdateProfile fetch error for user %d: %v", userID, err)
		return models.Profile{}, err
	}
	return saved.Profile, nil
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

// GetProfiles retrieves profiles, applying optional filters when provided.
func (s *ProfileService) GetProfiles(filters models.ProfileFilters) ([]models.UserProfile, error) {
	profiles, err := s.repo.GetAllWithFilters(filters)
	if err != nil {
		log.Printf("GetProfiles repository error: %v", err)
		return nil, err
	}
	return profiles, nil
}

// GetProfilesByUserIDs retrieves profiles indexed by user ID for the provided IDs.
func (s *ProfileService) GetProfilesByUserIDs(userIDs []int) (map[int]models.UserProfile, error) {
	profiles, err := s.repo.GetByUserIDs(userIDs)
	if err != nil {
		log.Printf("GetProfilesByUserIDs repository error: %v", err)
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

// EnqueueProfileSync schedules a profile synchronization attempt with the matching microservice.
func (s *ProfileService) EnqueueProfileSync(userID int) error {
	if s.profileOutboxRepo == nil {
		return errors.New("profile outbox repository not configured")
	}
	if err := s.profileOutboxRepo.Enqueue(userID); err != nil {
		log.Printf("EnqueueProfileSync enqueue error for user %d: %v", userID, err)
		return err
	}
	return nil
}

func parseVerificationToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, ErrInvalidVerificationToken
	}
	if !token.Valid {
		return nil, ErrInvalidVerificationToken
	}
	if !claimsVerified(claims) {
		return nil, ErrInvalidVerificationToken
	}
	return claims, nil
}

func claimsVerified(claims jwt.MapClaims) bool {
	if v, ok := claims["verified"]; ok {
		switch val := v.(type) {
		case bool:
			if val {
				return true
			}
		case string:
			if strings.EqualFold(val, "true") || strings.EqualFold(val, "verified") {
				return true
			}
		}
	}
	if v, ok := claims["status"].(string); ok {
		if strings.EqualFold(v, "verified") || strings.EqualFold(v, "approved") {
			return true
		}
	}
	if v, ok := claims["result"].(string); ok {
		if strings.EqualFold(v, "verified") || strings.EqualFold(v, "approved") {
			return true
		}
	}
	return false
}

func extractPhoneNumber(claims jwt.MapClaims) string {
	if v, ok := claims["phone_number"].(string); ok {
		return v
	}
	if v, ok := claims["phoneNumber"].(string); ok {
		return v
	}
	return ""
}

func getVerificationSecret(primaryEnv string) []byte {
	keys := []string{primaryEnv, "VERIFICATION_JWT_SECRET", "JWT_SECRET"}
	for _, key := range keys {
		if key == "" {
			continue
		}
		if secret := strings.TrimSpace(os.Getenv(key)); secret != "" {
			return []byte(secret)
		}
	}
	return []byte("secret")
}
