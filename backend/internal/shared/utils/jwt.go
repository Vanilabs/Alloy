package utils

import (
	"alloy/internal/shared/database/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SessionExpiry      = 7 * 24 * time.Hour  // 7 days
	RefreshTokenExpiry = 30 * 24 * time.Hour // 30 days
)

type JWTManager struct {
	JWTSecret          []byte
	RefreshTokenSecret []byte
	GoogleClientID     string
	FacebookAppID      string
	FacebookAppSecret  string
	AppleClientID      string
}

func NewJWTManager(jwtSecret, refreshTokenSecret string) *JWTManager {
	return &JWTManager{
		JWTSecret:          []byte(jwtSecret),
		RefreshTokenSecret: []byte(refreshTokenSecret),
	}
}

func (j *JWTManager) GenerateJWT(data *models.JWTData) (string, error) {
	tokenType := data.TokenType
	if tokenType == "" {
		tokenType = models.TokenTypeAccess
	}

	claims := models.Claims{
		UserID:    data.UserID,
		Email:     data.Email,
		Role:      data.Role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   data.UserID,
			ID:        data.TokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(SessionExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JWTSecret)
}

func (j *JWTManager) GenerateRefreshToken(userID string, tokenID string) (string, error) {
	claims := models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshTokenSecret)
}

func (j *JWTManager) ParseJWT(tokenStr string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func (j *JWTManager) ParseRefreshToken(tokenStr string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.RefreshTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func (j *JWTManager) VerifyRefreshToken(refreshToken string) (*models.Claims, error) {
	claims, err := j.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
