package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	MagicLinkTokenLength = 64

	MagicLinkTokenExpiry = 15 * time.Minute

	TokenTypeAccess TokenType = "access"
)

type JwtAuthData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type LoginResponse struct {
	Auth JwtAuthData `json:"auth"`
	User User        `json:"user"`
}

type JWTData struct {
	UserID    string
	Email     string
	TokenID   string
	Role      string
	TokenType TokenType
}

type Claims struct {
	UserID      string    `json:"sub"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	CountryCode string    `json:"country_code,omitempty"`
	Role        string    `json:"role"`
	TokenType   TokenType `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}
