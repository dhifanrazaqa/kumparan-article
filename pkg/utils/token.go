package utils

import (
	"time"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(user *models.User, accessTokenSecret string, refreshTokenSecret string) (string, string, error) {
	accessTokenClaims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := &jwt.RegisteredClaims{
		Subject: user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
