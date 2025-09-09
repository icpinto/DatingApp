package repositories

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/lib/pq"
)

func UpsertProfile(db *sql.DB, profile models.Profile) error {
	_, err := db.Exec(`
        INSERT INTO profiles (user_id, bio, gender, date_of_birth, location, interests)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (user_id)
        DO UPDATE SET bio = $2, gender = $3, date_of_birth = $4, location = $5, interests = $6, updated_at = NOW()`,
		profile.UserID, profile.Bio, profile.Gender, profile.DateOfBirth, profile.Location, pq.Array(profile.Interests))
	if err != nil {
		log.Printf("UpsertProfile exec error for user %d: %v", profile.UserID, err)
	}
	return err
}

func GetProfileByUserID(db *sql.DB, userID int) (models.Profile, error) {
	var profile models.Profile
	err := db.QueryRow(`
        SELECT id, user_id, bio, gender, date_of_birth, location, interests, created_at, updated_at
        FROM profiles WHERE user_id = $1`, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Bio, &profile.Gender,
		&profile.DateOfBirth, &profile.Location, pq.Array(&profile.Interests),
		&profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		log.Printf("GetProfileByUserID query error for user %d: %v", userID, err)
	}
	return profile, err
}

func GetAllProfiles(db *sql.DB) ([]models.Profile, error) {
	rows, err := db.Query(`
                SELECT id, user_id, bio, gender, date_of_birth, location, interests, created_at, updated_at
                FROM profiles`)
	if err != nil {
		log.Printf("GetAllProfiles query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(
			&profile.ID, &profile.UserID, &profile.Bio, &profile.Gender,
			&profile.DateOfBirth, &profile.Location, pq.Array(&profile.Interests),
			&profile.CreatedAt, &profile.UpdatedAt,
		); err != nil {
			log.Printf("GetAllProfiles scan error: %v", err)
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetAllProfiles rows error: %v", err)
		return nil, err
	}
	return profiles, nil
}
