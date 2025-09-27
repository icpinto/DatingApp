package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/icpinto/dating-app/models"
)

type corePreferencesDTO struct {
	UserID             int    `json:"user_id"`
	MinAge             int    `json:"min_age"`
	MaxAge             int    `json:"max_age"`
	Gender             string `json:"gender"`
	DrinkingHabit      string `json:"drinking_habit"`
	EducationLevel     string `json:"education_level"`
	SmokingHabit       string `json:"smoking_habit"`
	CountryOfResidence string `json:"country_of_residence"`
	OccupationStatus   string `json:"occupation_status"`
	CivilStatus        string `json:"civil_status"`
	Religion           string `json:"religion"`
	MinHeight          int    `json:"min_height"`
	MaxHeight          int    `json:"max_height"`
	FoodPreference     string `json:"food_preference"`
}

type profilePayload struct {
	UserID               int      `json:"user_id"`
	Bio                  string   `json:"bio"`
	Gender               string   `json:"gender"`
	DateOfBirth          string   `json:"date_of_birth"`
	LocationLegacy       string   `json:"location_legacy"`
	Interests            []string `json:"interests"`
	CivilStatus          string   `json:"civil_status"`
	Religion             string   `json:"religion"`
	ReligionDetail       string   `json:"religion_detail"`
	Caste                string   `json:"caste"`
	HeightCM             int      `json:"height_cm"`
	WeightKG             int      `json:"weight_kg"`
	DietaryPreference    string   `json:"dietary_preference"`
	Smoking              string   `json:"smoking"`
	Alcohol              string   `json:"alcohol"`
	Languages            []string `json:"languages"`
	PhoneNumber          string   `json:"phone_number"`
	ContactVerified      bool     `json:"contact_verified"`
	IdentityVerified     bool     `json:"identity_verified"`
	CountryCode          string   `json:"country_code"`
	Province             string   `json:"province"`
	District             string   `json:"district"`
	City                 string   `json:"city"`
	PostalCode           string   `json:"postal_code"`
	HighestEducation     string   `json:"highest_education"`
	FieldOfStudy         string   `json:"field_of_study"`
	Institution          string   `json:"institution"`
	EmploymentStatus     string   `json:"employment_status"`
	Occupation           string   `json:"occupation"`
	FatherOccupation     string   `json:"father_occupation"`
	MotherOccupation     string   `json:"mother_occupation"`
	SiblingsCount        int      `json:"siblings_count"`
	Siblings             string   `json:"siblings"`
	HoroscopeAvailable   bool     `json:"horoscope_available"`
	BirthTime            string   `json:"birth_time"`
	BirthPlace           string   `json:"birth_place"`
	SinhalaRaasi         string   `json:"sinhala_raasi"`
	Nakshatra            string   `json:"nakshatra"`
	Horoscope            string   `json:"horoscope"`
	ProfileImageURL      string   `json:"profile_image_url"`
	ProfileImageThumbURL string   `json:"profile_image_thumb_url"`
	Verified             bool     `json:"verified"`
	ModerationStatus     string   `json:"moderation_status"`
	LastActiveAt         string   `json:"last_active_at"`
	Metadata             string   `json:"metadata"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
}

type profileDTO struct {
	ID int `json:"id"`
	profilePayload
}

func newCorePreferencesDTO(prefs models.CorePreferences) corePreferencesDTO {
	return corePreferencesDTO{
		UserID:             prefs.UserID,
		MinAge:             prefs.MinAge,
		MaxAge:             prefs.MaxAge,
		Gender:             prefs.Gender,
		DrinkingHabit:      prefs.DrinkingHabit,
		EducationLevel:     prefs.EducationLevel,
		SmokingHabit:       prefs.SmokingHabit,
		CountryOfResidence: prefs.CountryOfResidence,
		OccupationStatus:   prefs.OccupationStatus,
		CivilStatus:        prefs.CivilStatus,
		Religion:           prefs.Religion,
		MinHeight:          prefs.MinHeight,
		MaxHeight:          prefs.MaxHeight,
		FoodPreference:     prefs.FoodPreference,
	}
}

func (d corePreferencesDTO) toModel() models.CorePreferences {
	return models.CorePreferences{
		UserID:             d.UserID,
		MinAge:             d.MinAge,
		MaxAge:             d.MaxAge,
		Gender:             d.Gender,
		DrinkingHabit:      d.DrinkingHabit,
		EducationLevel:     d.EducationLevel,
		SmokingHabit:       d.SmokingHabit,
		CountryOfResidence: d.CountryOfResidence,
		OccupationStatus:   d.OccupationStatus,
		CivilStatus:        d.CivilStatus,
		Religion:           d.Religion,
		MinHeight:          d.MinHeight,
		MaxHeight:          d.MaxHeight,
		FoodPreference:     d.FoodPreference,
	}
}

func newProfilePayload(profile models.Profile) profilePayload {
	return profilePayload{
		UserID:               profile.UserID,
		Bio:                  profile.Bio,
		Gender:               profile.Gender,
		DateOfBirth:          profile.DateOfBirth,
		LocationLegacy:       profile.LocationLegacy,
		Interests:            profile.Interests,
		CivilStatus:          profile.CivilStatus,
		Religion:             profile.Religion,
		ReligionDetail:       profile.ReligionDetail,
		Caste:                profile.Caste,
		HeightCM:             profile.HeightCM,
		WeightKG:             profile.WeightKG,
		DietaryPreference:    profile.DietaryPreference,
		Smoking:              profile.Smoking,
		Alcohol:              profile.Alcohol,
		Languages:            profile.Languages,
		PhoneNumber:          profile.PhoneNumber,
		ContactVerified:      profile.ContactVerified,
		IdentityVerified:     profile.IdentityVerified,
		CountryCode:          profile.CountryCode,
		Province:             profile.Province,
		District:             profile.District,
		City:                 profile.City,
		PostalCode:           profile.PostalCode,
		HighestEducation:     profile.HighestEducation,
		FieldOfStudy:         profile.FieldOfStudy,
		Institution:          profile.Institution,
		EmploymentStatus:     profile.EmploymentStatus,
		Occupation:           profile.Occupation,
		FatherOccupation:     profile.FatherOccupation,
		MotherOccupation:     profile.MotherOccupation,
		SiblingsCount:        profile.SiblingsCount,
		Siblings:             profile.Siblings,
		HoroscopeAvailable:   profile.HoroscopeAvailable,
		BirthTime:            profile.BirthTime,
		BirthPlace:           profile.BirthPlace,
		SinhalaRaasi:         profile.SinhalaRaasi,
		Nakshatra:            profile.Nakshatra,
		Horoscope:            profile.Horoscope,
		ProfileImageURL:      profile.ProfileImageURL,
		ProfileImageThumbURL: profile.ProfileImageThumbURL,
		Verified:             profile.Verified,
		ModerationStatus:     profile.ModerationStatus,
		LastActiveAt:         profile.LastActiveAt,
		Metadata:             profile.Metadata,
		CreatedAt:            profile.CreatedAt,
		UpdatedAt:            profile.UpdatedAt,
	}
}

func (d profileDTO) toModel() models.Profile {
	payload := d.profilePayload
	return models.Profile{
		ID:                   d.ID,
		UserID:               payload.UserID,
		Bio:                  payload.Bio,
		Gender:               payload.Gender,
		DateOfBirth:          payload.DateOfBirth,
		LocationLegacy:       payload.LocationLegacy,
		Interests:            payload.Interests,
		CivilStatus:          payload.CivilStatus,
		Religion:             payload.Religion,
		ReligionDetail:       payload.ReligionDetail,
		Caste:                payload.Caste,
		HeightCM:             payload.HeightCM,
		WeightKG:             payload.WeightKG,
		DietaryPreference:    payload.DietaryPreference,
		Smoking:              payload.Smoking,
		Alcohol:              payload.Alcohol,
		Languages:            payload.Languages,
		PhoneNumber:          payload.PhoneNumber,
		ContactVerified:      payload.ContactVerified,
		IdentityVerified:     payload.IdentityVerified,
		CountryCode:          payload.CountryCode,
		Province:             payload.Province,
		District:             payload.District,
		City:                 payload.City,
		PostalCode:           payload.PostalCode,
		HighestEducation:     payload.HighestEducation,
		FieldOfStudy:         payload.FieldOfStudy,
		Institution:          payload.Institution,
		EmploymentStatus:     payload.EmploymentStatus,
		Occupation:           payload.Occupation,
		FatherOccupation:     payload.FatherOccupation,
		MotherOccupation:     payload.MotherOccupation,
		SiblingsCount:        payload.SiblingsCount,
		Siblings:             payload.Siblings,
		HoroscopeAvailable:   payload.HoroscopeAvailable,
		BirthTime:            payload.BirthTime,
		BirthPlace:           payload.BirthPlace,
		SinhalaRaasi:         payload.SinhalaRaasi,
		Nakshatra:            payload.Nakshatra,
		Horoscope:            payload.Horoscope,
		ProfileImageURL:      payload.ProfileImageURL,
		ProfileImageThumbURL: payload.ProfileImageThumbURL,
		Verified:             payload.Verified,
		ModerationStatus:     payload.ModerationStatus,
		LastActiveAt:         payload.LastActiveAt,
		Metadata:             payload.Metadata,
		CreatedAt:            payload.CreatedAt,
		UpdatedAt:            payload.UpdatedAt,
	}
}

// MatchService communicates with the external matching microservice.
type MatchService struct {
	client  *http.Client
	baseURL string
}

// NewMatchService creates a MatchService with the provided base URL.
func NewMatchService(baseURL string) *MatchService {
	if baseURL == "" {
		baseURL = "http://localhost:8003"
	}
	return &MatchService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// GetMatches fetches match candidates for a user from the microservice.
func (s *MatchService) GetMatches(ctx context.Context, userID int, rawQuery string) ([]models.MatchCandidate, error) {
	endpoint := fmt.Sprintf("%s/matches/%d", s.baseURL, userID)
	if rawQuery != "" {
		endpoint = endpoint + "?" + rawQuery
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("match service returned status %d", resp.StatusCode)
	}

	var matches []models.MatchCandidate
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		return nil, err
	}
	return matches, nil
}

// SaveCorePreferences sends a request to create the user's core preferences in the matching microservice.
func (s *MatchService) SaveCorePreferences(ctx context.Context, prefs models.CorePreferences) (models.CorePreferences, error) {
	return s.sendCorePreferences(ctx, http.MethodPost, prefs)
}

// UpdateCorePreferences sends a request to update the user's core preferences in the matching microservice.
func (s *MatchService) UpdateCorePreferences(ctx context.Context, prefs models.CorePreferences) (models.CorePreferences, error) {
	return s.sendCorePreferences(ctx, http.MethodPut, prefs)
}

// GetCorePreferences retrieves the core preferences for the given user from the matching microservice.
func (s *MatchService) GetCorePreferences(ctx context.Context, userID int) (models.CorePreferences, error) {
	endpoint := fmt.Sprintf("%s/core-preferences/%d", s.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return models.CorePreferences{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return models.CorePreferences{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return models.CorePreferences{}, fmt.Errorf("core preferences not found")
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return models.CorePreferences{}, fmt.Errorf("match service returned status %d", resp.StatusCode)
	}

	var dto corePreferencesDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return models.CorePreferences{}, err
	}

	return dto.toModel(), nil
}

func (s *MatchService) sendCorePreferences(ctx context.Context, method string, prefs models.CorePreferences) (models.CorePreferences, error) {
	endpoint := fmt.Sprintf("%s/core-preferences", s.baseURL)

	payload, err := json.Marshal(newCorePreferencesDTO(prefs))
	if err != nil {
		return models.CorePreferences{}, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(payload))
	if err != nil {
		return models.CorePreferences{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return models.CorePreferences{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return models.CorePreferences{}, fmt.Errorf("match service returned status %d", resp.StatusCode)
	}

	var saved corePreferencesDTO
	if err := json.NewDecoder(resp.Body).Decode(&saved); err != nil {
		return models.CorePreferences{}, err
	}

	return saved.toModel(), nil
}

// UpsertProfile sends the given profile to the matching microservice.
func (s *MatchService) UpsertProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	endpoint := fmt.Sprintf("%s/profiles", s.baseURL)
	method := http.MethodPost
	if profile.ID != 0 {
		method = http.MethodPut
	}

	payload, err := json.Marshal(newProfilePayload(profile))
	if err != nil {
		return models.Profile{}, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(payload))
	if err != nil {
		return models.Profile{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return models.Profile{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return models.Profile{}, fmt.Errorf("match service returned status %d", resp.StatusCode)
	}

	var saved profileDTO
	if err := json.NewDecoder(resp.Body).Decode(&saved); err != nil {
		return models.Profile{}, err
	}

	return saved.toModel(), nil
}
