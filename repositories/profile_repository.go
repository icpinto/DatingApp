package repositories

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/lib/pq"
)

// ProfileRepository handles CRUD operations for profiles.
type ProfileRepository struct {
	db *sql.DB
}

// NewProfileRepository creates a new ProfileRepository.
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Upsert creates or updates a profile record.
func (r *ProfileRepository) Upsert(profile models.Profile) error {
	_, err := r.db.Exec(`
        INSERT INTO profiles (user_id, bio, gender, date_of_birth, location, interests, profile_image)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (user_id)
        DO UPDATE SET bio = EXCLUDED.bio, gender = EXCLUDED.gender, date_of_birth = EXCLUDED.date_of_birth, location = EXCLUDED.location, interests = EXCLUDED.interests, profile_image = CASE WHEN EXCLUDED.profile_image <> '' THEN EXCLUDED.profile_image ELSE profiles.profile_image END, updated_at = NOW()`,
		profile.UserID, profile.Bio, profile.Gender, profile.DateOfBirth, profile.Location, pq.Array(profile.Interests), profile.ProfileImage)
	if err != nil {
		log.Printf("ProfileRepository.Upsert error for user %d: %v", profile.UserID, err)
	}
	return err
}

// GetByUserID retrieves a profile for the specified user ID.
func (r *ProfileRepository) GetByUserID(userID int) (models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.QueryRow(`
       SELECT p.id, p.user_id, u.username, p.bio, p.gender, p.date_of_birth, p.location, p.interests, p.profile_image, p.created_at, p.updated_at
       FROM profiles p JOIN users u ON p.user_id = u.id WHERE p.user_id = $1`, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Username, &profile.Bio, &profile.Gender,
		&profile.DateOfBirth, &profile.Location, pq.Array(&profile.Interests), &profile.ProfileImage,
		&profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		log.Printf("ProfileRepository.GetByUserID query error for user %d: %v", userID, err)
	}
	return profile, err
}

// GetAll retrieves all profiles.
func (r *ProfileRepository) GetAll() ([]models.UserProfile, error) {
	rows, err := r.db.Query(`
               SELECT p.id, p.user_id, u.username, p.bio, p.gender, p.date_of_birth, p.location, p.interests, p.profile_image, p.created_at, p.updated_at
               FROM profiles p JOIN users u ON p.user_id = u.id`)
	if err != nil {
		log.Printf("ProfileRepository.GetAll query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var profiles []models.UserProfile
	for rows.Next() {
		var profile models.UserProfile
		if err := rows.Scan(
			&profile.ID, &profile.UserID, &profile.Username, &profile.Bio, &profile.Gender,
			&profile.DateOfBirth, &profile.Location, pq.Array(&profile.Interests), &profile.ProfileImage,
			&profile.CreatedAt, &profile.UpdatedAt,
		); err != nil {
			log.Printf("ProfileRepository.GetAll scan error: %v", err)
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ProfileRepository.GetAll rows error: %v", err)
		return nil, err
	}
	return profiles, nil
}

// Delete removes a profile by user ID.
func (r *ProfileRepository) Delete(userID int) error {
	if _, err := r.db.Exec(`DELETE FROM profiles WHERE user_id = $1`, userID); err != nil {
		log.Printf("ProfileRepository.Delete error for user %d: %v", userID, err)
		return err
	}
	return nil
}
