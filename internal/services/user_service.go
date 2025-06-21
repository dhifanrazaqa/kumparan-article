package services

import (
	"context"
	"errors"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/dhifanrazaqa/kumparan-article/internal/repositories"
	"github.com/dhifanrazaqa/kumparan-article/pkg/utils"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var (
	ErrForbidden = errors.New("you do not have permission to perform this action")
)

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
	GetUsers(ctx context.Context) ([]models.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest, currentUserID string) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id string, currentUserID string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	_, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := &models.User{
		Username:       req.Username,
		Name:           req.Name,
		HashedPassword: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetUsers(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	return responses, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest, currentUserID string) (*models.UserResponse, error) {
	if id != currentUserID {
		return nil, ErrForbidden
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("failed to process new password")
		}
		user.HashedPassword = hashedPassword
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string, currentUserID string) error {
	if id != currentUserID {
		return ErrForbidden
	}
	return s.userRepo.Delete(ctx, id)
}
