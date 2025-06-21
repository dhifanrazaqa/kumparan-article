package services

import (
	"context"
	"errors"
	"testing"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/dhifanrazaqa/kumparan-article/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil {
		user.ID = "1e9e25d9-5b6e-4dbc-a5fd-a0bd8aa209de"
	}
	return args.Error(0)
}

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}
func (m *MockUserRepo) FindAll(ctx context.Context) ([]models.User, error)  { return nil, nil }
func (m *MockUserRepo) Update(ctx context.Context, user *models.User) error { return nil }
func (m *MockUserRepo) Delete(ctx context.Context, id string) error         { return nil }

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)

	userService := NewUserService(mockRepo)

	t.Run("sukses membuat pengguna baru", func(t *testing.T) {
		mockRepo.On("FindByUsername", mock.Anything, "newuser").Return(nil, repositories.ErrUserNotFound).Once()
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()

		req := models.CreateUserRequest{Username: "newuser", Password: "password123"}
		user, err := userService.CreateUser(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "newuser", user.Username)
		assert.Equal(t, "1e9e25d9-5b6e-4dbc-a5fd-a0bd8aa209de", user.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("gagal karena username sudah terdaftar", func(t *testing.T) {
		existingUser := &models.User{ID: "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6", Username: "existinguser"}
		mockRepo.On("FindByUsername", mock.Anything, "existinguser").Return(existingUser, nil).Once()

		req := models.CreateUserRequest{Username: "existinguser", Password: "password123"}
		_, err := userService.CreateUser(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, ErrUserAlreadyExists, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("gagal karena error di repository saat create", func(t *testing.T) {
		dbError := errors.New("koneksi database terputus")
		mockRepo.On("FindByUsername", mock.Anything, "anotheruser").Return(nil, repositories.ErrUserNotFound).Once()
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(dbError).Once()

		req := models.CreateUserRequest{Username: "anotheruser", Password: "password123"}
		_, err := userService.CreateUser(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)

		mockRepo.AssertExpectations(t)
	})
}
