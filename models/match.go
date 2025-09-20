package models

import "encoding/json"

// MatchCandidate represents the response from the matching microservice.
type MatchCandidate struct {
	UserID  int             `json:"user_id"`
	Score   float64         `json:"score"`
	Reasons json.RawMessage `json:"reasons" swaggertype:"object"`
}

// MatchedProfile combines a user profile with a compatibility score.
type MatchedProfile struct {
	UserProfile
	Score   float64         `json:"score"`
	Reasons json.RawMessage `json:"reasons" swaggertype:"object"`
}
