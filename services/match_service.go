package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/icpinto/dating-app/models"
)

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
