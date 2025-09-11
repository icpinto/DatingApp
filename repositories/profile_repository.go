package repositories

import (
	"database/sql"
	"encoding/json"
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
	// Ensure metadata is valid JSON; default to empty object if not provided or invalid
	metadata := []byte("{}")
	if profile.Metadata != "" {
		if json.Valid([]byte(profile.Metadata)) {
			metadata = []byte(profile.Metadata)
		} else {
			log.Printf("ProfileRepository.Upsert invalid metadata for user %d: %s", profile.UserID, profile.Metadata)
		}
	}

	civilStatus := sql.NullString{String: profile.CivilStatus, Valid: profile.CivilStatus != ""}
	dietaryPreference := sql.NullString{String: profile.DietaryPreference, Valid: profile.DietaryPreference != ""}
	smoking := sql.NullString{String: profile.Smoking, Valid: profile.Smoking != ""}
	alcohol := sql.NullString{String: profile.Alcohol, Valid: profile.Alcohol != ""}
	highestEducation := sql.NullString{String: profile.HighestEducation, Valid: profile.HighestEducation != ""}
	employmentStatus := sql.NullString{String: profile.EmploymentStatus, Valid: profile.EmploymentStatus != ""}

	_, err := r.db.Exec(`
INSERT INTO profiles (
user_id, bio, gender, date_of_birth, location_legacy, interests, civil_status, religion, religion_detail,
caste, height_cm, weight_kg, dietary_preference, smoking, alcohol, languages,
country_code, province, district, city, postal_code,
highest_education, field_of_study, institution, employment_status, occupation,
father_occupation, mother_occupation, siblings_count, siblings,
horoscope_available, birth_time, birth_place, sinhala_raasi, nakshatra, horoscope,
profile_image_url, profile_image_thumb_url, verified, moderation_status, last_active_at, metadata)
VALUES (
$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
$31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42)
ON CONFLICT (user_id)
DO UPDATE SET bio = EXCLUDED.bio, gender = EXCLUDED.gender, date_of_birth = EXCLUDED.date_of_birth,
location_legacy = EXCLUDED.location_legacy, interests = EXCLUDED.interests, civil_status = EXCLUDED.civil_status,
religion = EXCLUDED.religion, religion_detail = EXCLUDED.religion_detail, caste = EXCLUDED.caste,
height_cm = EXCLUDED.height_cm, weight_kg = EXCLUDED.weight_kg,
dietary_preference = EXCLUDED.dietary_preference, smoking = EXCLUDED.smoking, alcohol = EXCLUDED.alcohol,
languages = EXCLUDED.languages, country_code = EXCLUDED.country_code, province = EXCLUDED.province,
district = EXCLUDED.district, city = EXCLUDED.city, postal_code = EXCLUDED.postal_code,
highest_education = EXCLUDED.highest_education, field_of_study = EXCLUDED.field_of_study,
institution = EXCLUDED.institution, employment_status = EXCLUDED.employment_status,
occupation = EXCLUDED.occupation, father_occupation = EXCLUDED.father_occupation,
mother_occupation = EXCLUDED.mother_occupation, siblings_count = EXCLUDED.siblings_count,
siblings = EXCLUDED.siblings, horoscope_available = EXCLUDED.horoscope_available,
birth_time = EXCLUDED.birth_time, birth_place = EXCLUDED.birth_place,
sinhala_raasi = EXCLUDED.sinhala_raasi, nakshatra = EXCLUDED.nakshatra,
horoscope = EXCLUDED.horoscope,
profile_image_url = CASE WHEN EXCLUDED.profile_image_url <> '' THEN EXCLUDED.profile_image_url ELSE profiles.profile_image_url END,
profile_image_thumb_url = CASE WHEN EXCLUDED.profile_image_thumb_url <> '' THEN EXCLUDED.profile_image_thumb_url ELSE profiles.profile_image_thumb_url END,
verified = EXCLUDED.verified, moderation_status = EXCLUDED.moderation_status,
last_active_at = EXCLUDED.last_active_at, metadata = EXCLUDED.metadata,
updated_at = NOW()`,
		profile.UserID, profile.Bio, profile.Gender, profile.DateOfBirth, profile.LocationLegacy,
		pq.Array(profile.Interests), civilStatus, profile.Religion, profile.ReligionDetail,
		profile.Caste, profile.HeightCM, profile.WeightKG, dietaryPreference, smoking, alcohol,
		pq.Array(profile.Languages), profile.CountryCode, profile.Province, profile.District, profile.City, profile.PostalCode,
		highestEducation, profile.FieldOfStudy, profile.Institution, employmentStatus, profile.Occupation,
		profile.FatherOccupation, profile.MotherOccupation, profile.SiblingsCount, profile.Siblings,
		profile.HoroscopeAvailable, profile.BirthTime, profile.BirthPlace, profile.SinhalaRaasi, profile.Nakshatra, profile.Horoscope,
		profile.ProfileImageURL, profile.ProfileImageThumbURL, profile.Verified, profile.ModerationStatus,
		profile.LastActiveAt, metadata)
	if err != nil {
		log.Printf("ProfileRepository.Upsert error for user %d: %v", profile.UserID, err)
	}
	return err
}

// GetByUserID retrieves a profile for the specified user ID.
func (r *ProfileRepository) GetByUserID(userID int) (models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.QueryRow(`
       SELECT p.id, p.user_id, u.username, p.bio, p.gender, p.date_of_birth,
              COALESCE(p.location_legacy, ''), COALESCE(p.interests, ARRAY[]::text[]),
              COALESCE(p.civil_status::text, ''), COALESCE(p.religion, ''), COALESCE(p.religion_detail, ''), COALESCE(p.caste, ''),
              COALESCE(p.height_cm, 0), COALESCE(p.weight_kg, 0), COALESCE(p.dietary_preference::text, ''), COALESCE(p.smoking::text, ''), COALESCE(p.alcohol::text, ''),
              COALESCE(p.languages, ARRAY[]::text[]),
              COALESCE(p.country_code, ''), COALESCE(p.province, ''), COALESCE(p.district, ''), COALESCE(p.city, ''), COALESCE(p.postal_code, ''),
              COALESCE(p.highest_education::text, ''), COALESCE(p.field_of_study, ''), COALESCE(p.institution, ''), COALESCE(p.employment_status::text, ''), COALESCE(p.occupation, ''),
              COALESCE(p.father_occupation, ''), COALESCE(p.mother_occupation, ''), COALESCE(p.siblings_count, 0), COALESCE(p.siblings::text, ''),
              COALESCE(p.horoscope_available, false), COALESCE(p.birth_time::text, ''), COALESCE(p.birth_place, ''), COALESCE(p.sinhala_raasi, ''), COALESCE(p.nakshatra, ''), COALESCE(p.horoscope::text, ''),
              COALESCE(p.profile_image_url, ''), COALESCE(p.profile_image_thumb_url, ''), COALESCE(p.verified, false), COALESCE(p.moderation_status, ''), COALESCE(p.last_active_at::text, ''), COALESCE(p.metadata::text, ''),
              p.created_at, p.updated_at
       FROM profiles p JOIN users u ON p.user_id = u.id WHERE p.user_id = $1`, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Username, &profile.Bio, &profile.Gender,
		&profile.DateOfBirth, &profile.LocationLegacy, pq.Array(&profile.Interests),
		&profile.CivilStatus, &profile.Religion, &profile.ReligionDetail, &profile.Caste,
		&profile.HeightCM, &profile.WeightKG, &profile.DietaryPreference, &profile.Smoking, &profile.Alcohol,
		pq.Array(&profile.Languages), &profile.CountryCode, &profile.Province,
		&profile.District, &profile.City, &profile.PostalCode,
		&profile.HighestEducation, &profile.FieldOfStudy, &profile.Institution, &profile.EmploymentStatus, &profile.Occupation,
		&profile.FatherOccupation, &profile.MotherOccupation, &profile.SiblingsCount, &profile.Siblings,
		&profile.HoroscopeAvailable, &profile.BirthTime, &profile.BirthPlace, &profile.SinhalaRaasi, &profile.Nakshatra, &profile.Horoscope,
		&profile.ProfileImageURL, &profile.ProfileImageThumbURL, &profile.Verified, &profile.ModerationStatus, &profile.LastActiveAt, &profile.Metadata,
		&profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		log.Printf("ProfileRepository.GetByUserID query error for user %d: %v", userID, err)
	}
	return profile, err
}

// GetAll retrieves all profiles.
func (r *ProfileRepository) GetAll() ([]models.UserProfile, error) {
	rows, err := r.db.Query(`
               SELECT p.id, p.user_id, u.username, p.bio, p.gender, p.date_of_birth,
                      COALESCE(p.location_legacy, ''), COALESCE(p.interests, ARRAY[]::text[]),
                      COALESCE(p.civil_status::text, ''), COALESCE(p.religion, ''), COALESCE(p.religion_detail, ''), COALESCE(p.caste, ''),
                      COALESCE(p.height_cm, 0), COALESCE(p.weight_kg, 0), COALESCE(p.dietary_preference::text, ''), COALESCE(p.smoking::text, ''), COALESCE(p.alcohol::text, ''),
                      COALESCE(p.languages, ARRAY[]::text[]),
                      COALESCE(p.country_code, ''), COALESCE(p.province, ''), COALESCE(p.district, ''), COALESCE(p.city, ''), COALESCE(p.postal_code, ''),
                      COALESCE(p.highest_education::text, ''), COALESCE(p.field_of_study, ''), COALESCE(p.institution, ''), COALESCE(p.employment_status::text, ''), COALESCE(p.occupation, ''),
                      COALESCE(p.father_occupation, ''), COALESCE(p.mother_occupation, ''), COALESCE(p.siblings_count, 0), COALESCE(p.siblings::text, ''),
                      COALESCE(p.horoscope_available, false), COALESCE(p.birth_time::text, ''), COALESCE(p.birth_place, ''), COALESCE(p.sinhala_raasi, ''), COALESCE(p.nakshatra, ''), COALESCE(p.horoscope::text, ''),
                      COALESCE(p.profile_image_url, ''), COALESCE(p.profile_image_thumb_url, ''), COALESCE(p.verified, false), COALESCE(p.moderation_status, ''), COALESCE(p.last_active_at::text, ''), COALESCE(p.metadata::text, ''),
                      p.created_at, p.updated_at
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
			&profile.DateOfBirth, &profile.LocationLegacy, pq.Array(&profile.Interests),
			&profile.CivilStatus, &profile.Religion, &profile.ReligionDetail, &profile.Caste,
			&profile.HeightCM, &profile.WeightKG, &profile.DietaryPreference, &profile.Smoking, &profile.Alcohol,
			pq.Array(&profile.Languages), &profile.CountryCode, &profile.Province,
			&profile.District, &profile.City, &profile.PostalCode,
			&profile.HighestEducation, &profile.FieldOfStudy, &profile.Institution, &profile.EmploymentStatus, &profile.Occupation,
			&profile.FatherOccupation, &profile.MotherOccupation, &profile.SiblingsCount, &profile.Siblings,
			&profile.HoroscopeAvailable, &profile.BirthTime, &profile.BirthPlace, &profile.SinhalaRaasi, &profile.Nakshatra, &profile.Horoscope,
			&profile.ProfileImageURL, &profile.ProfileImageThumbURL, &profile.Verified, &profile.ModerationStatus, &profile.LastActiveAt, &profile.Metadata,
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

// getEnumValues returns the labels for a given PostgreSQL enum type.
func (r *ProfileRepository) getEnumValues(enumType string) ([]string, error) {
	rows, err := r.db.Query(`SELECT enumlabel FROM pg_enum WHERE enumtypid = (SELECT oid FROM pg_type WHERE typname = $1)`, enumType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, rows.Err()
}

// GetProfileEnums fetches enum values for profile-related fields.
func (r *ProfileRepository) GetProfileEnums() (models.ProfileEnums, error) {
	enums := models.ProfileEnums{}
	var err error
	if enums.CivilStatus, err = r.getEnumValues("civil_status_type"); err != nil {
		return enums, err
	}
	if enums.DietaryPreference, err = r.getEnumValues("dietary_pref_type"); err != nil {
		return enums, err
	}
	if enums.HabitFrequency, err = r.getEnumValues("habit_freq_type"); err != nil {
		return enums, err
	}
	if enums.EducationLevel, err = r.getEnumValues("education_level_type"); err != nil {
		return enums, err
	}
	if enums.EmploymentStatus, err = r.getEnumValues("employment_status_type"); err != nil {
		return enums, err
	}
	return enums, nil
}
