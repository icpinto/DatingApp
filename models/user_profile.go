package models

import "github.com/golang-jwt/jwt"

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}

// Claims defines the structure of the JWT payload
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Profile struct {
	ID                   int      `json:"id"`
	UserID               int      `json:"user_id"` // Foreign key to users table
	Bio                  string   `json:"bio"`
	Gender               string   `json:"gender"`
	DateOfBirth          string   `json:"date_of_birth"`
	LocationLegacy       string   `json:"location"`
	Interests            []string `json:"interests"` // Array of interests
	Languages            []string `json:"languages"`
	CountryCode          string   `json:"country_code"`
	Province             string   `json:"province"`
	District             string   `json:"district"`
	City                 string   `json:"city"`
	PostalCode           string   `json:"postal_code"`
	ProfileImageURL      string   `json:"profile_image_url"`
	ProfileImageThumbURL string   `json:"profile_image_thumb_url"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
}

type UserProfile struct {
	Profile
	Username string `json:"username"`
}
