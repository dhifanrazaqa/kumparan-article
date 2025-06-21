package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/dhifanrazaqa/kumparan-article/internal/repositories"
	"github.com/dhifanrazaqa/kumparan-article/pkg/utils"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

type AuthService interface {
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*models.AuthResponse, error)
}

type authService struct {
	userRepo           repositories.UserRepository
	jwtSecret          string
	refreshTokenSecret string
	redisClient        *redis.Client
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret, refreshTokenSecret string, redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		refreshTokenSecret: refreshTokenSecret,
		redisClient:        redisClient,
	}
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(req.Password, user.HashedPassword) {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user, s.jwtSecret, s.refreshTokenSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	err = s.redisClient.Set(refreshToken, user.ID, 7*24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to save refresh token to Redis: %v", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshTokenString string) (*models.AuthResponse, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidRefreshToken
	}

	userID, err := s.redisClient.Get(refreshTokenString).Result()
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	s.redisClient.Del(refreshTokenString)

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	newAccessToken, newRefreshToken, err := utils.GenerateTokens(user, s.jwtSecret, s.refreshTokenSecret)
	if err != nil {
		return nil, errors.New("failed to generate new token")
	}

	err = s.redisClient.Set(newRefreshToken, user.ID, 7*24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to save new refresh token to Redis: %v", err)
	}

	return &models.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}