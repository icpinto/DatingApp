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
	UserID             int    `json:"userId"`
	MinAge             int    `json:"minAge"`
	MaxAge             int    `json:"maxAge"`
	Gender             string `json:"gender"`
	DrinkingHabit      string `json:"drinkingHabit"`
	EducationLevel     string `json:"educationLevel"`
	SmokingHabit       string `json:"smokingHabit"`
	CountryOfResidence string `json:"countryOfResidence"`
	OccupationStatus   string `json:"occupationStatus"`
	CivilStatus        string `json:"civilStatus"`
	Religion           string `json:"religion"`
	MinHeight          int    `json:"minHeight"`
	MaxHeight          int    `json:"maxHeight"`
	FoodPreference     string `json:"foodPreference"`
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
