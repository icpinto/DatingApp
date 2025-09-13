package models

import "github.com/golang-jwt/jwt"

// Claims defines the structure of the JWT payload
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
