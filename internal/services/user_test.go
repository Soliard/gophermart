// internal/services/user_test.go
package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/repository"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	users            map[string]*models.User
	shouldFailExist  bool
	shouldFailCreate bool
	existsResult     bool
}

func (m *mockUserRepository) UserExists(ctx context.Context, login string) (bool, error) {
	if m.shouldFailExist {
		return false, errors.New("db error")
	}

	if m.existsResult {
		return true, nil // принудительно возвращаем что пользователь существует
	}

	_, exists := m.users[login]
	return exists, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.shouldFailCreate {
		return errors.New("create error")
	}

	if m.users == nil {
		m.users = make(map[string]*models.User)
	}
	m.users[user.Login] = user
	return nil
}

func (m *mockUserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	user, exists := m.users[login]
	if !exists {
		return nil, errs.UserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) UpdateLoginTime(ctx context.Context, userID string, time time.Time) error {
	return nil
}

type mockJWTService struct {
	shouldFail    bool
	tokenToReturn string
}

func (m *mockJWTService) GenerateToken(u *models.User) (string, error) {
	if m.shouldFail {
		return "", errors.New("jwt error")
	}

	if m.tokenToReturn != "" {
		return m.tokenToReturn, nil
	}

	return "fake-jwt-token", nil
}

func (m *mockJWTService) GetClaims(token string) (*UserContext, error) {
	return &UserContext{ID: "user-123"}, nil
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name        string
		request     *dto.RegisterRequest
		setupMocks  func() (repository.UserRepositoryInterface, JWTServiceInterface)
		wantErr     error
		checkResult func(*testing.T, *models.User)
	}{
		{
			name: "successful registration",
			request: &dto.RegisterRequest{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() (repository.UserRepositoryInterface, JWTServiceInterface) {
				return &mockUserRepository{}, &mockJWTService{}
			},
			wantErr: nil,
			checkResult: func(t *testing.T, user *models.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "testuser", user.Login)
				assert.NotEmpty(t, user.ID)
				assert.NotEmpty(t, user.PasswordHash)
			},
		},
		{
			name: "login already exists",
			request: &dto.RegisterRequest{
				Login:    "existinguser",
				Password: "password123",
			},
			setupMocks: func() (repository.UserRepositoryInterface, JWTServiceInterface) {
				return &mockUserRepository{existsResult: true}, &mockJWTService{}
			},
			wantErr: errs.LoginAlreadyExists,
			checkResult: func(t *testing.T, user *models.User) {
				assert.Nil(t, user)
			},
		},
		{
			name: "empty login",
			request: &dto.RegisterRequest{
				Login:    "",
				Password: "password123",
			},
			setupMocks: func() (repository.UserRepositoryInterface, JWTServiceInterface) {
				return &mockUserRepository{}, &mockJWTService{}
			},
			wantErr: errs.EmptyLoginOrPassword,
			checkResult: func(t *testing.T, user *models.User) {
				assert.Nil(t, user)
			},
		},
		{
			name: "repository error on UserExists",
			request: &dto.RegisterRequest{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() (repository.UserRepositoryInterface, JWTServiceInterface) {
				return &mockUserRepository{shouldFailExist: true}, &mockJWTService{}
			},
			wantErr: errors.New("db error"),
			checkResult: func(t *testing.T, user *models.User) {
				assert.Nil(t, user)
			},
		},
		{
			name: "repository error on Create",
			request: &dto.RegisterRequest{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() (repository.UserRepositoryInterface, JWTServiceInterface) {
				return &mockUserRepository{shouldFailCreate: true}, &mockJWTService{}
			},
			wantErr: errors.New("create error"),
			checkResult: func(t *testing.T, user *models.User) {
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			mockRepo, mockJWT := tt.setupMocks()

			// Создаем сервис
			service := NewUserService(mockRepo, mockJWT)

			// Выполняем тест
			user, err := service.Register(context.Background(), tt.request)

			// Проверяем ошибку
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем результат
			tt.checkResult(t, user)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name       string
		request    *dto.LoginRequest
		setupMocks func() (repository.UserRepositoryInterface, JWTServiceInterface)
		wantErr    error
		wantToken  string
	}{
		{
			name: "successful login",
			request: &dto.LoginRequest{
				Login:    "testuser",
				Password: "password123",
			},
			setupMocks: func() (repository.UserRepositoryInterface, JWTServiceInterface) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &models.User{
					ID:           "user-123",
					Login:        "testuser",
					PasswordHash: string(hashedPassword),
				}

				mockRepo := &mockUserRepository{
					users: map[string]*models.User{"testuser": user},
				}
				mockJWT := &mockJWTService{tokenToReturn: "success-token"}

				return mockRepo, mockJWT
			},
			wantErr:   nil,
			wantToken: "success-token",
		},
		{
			name: "user not found",
			request: &dto.LoginRequest{
				Login:    "nonexistent",
				Password: "password123",
			},
			setupMocks: func() (repository.UserRepositoryInterface, JWTServiceInterface) {
				return &mockUserRepository{}, &mockJWTService{}
			},
			wantErr:   errs.WrongLoginOrPassword,
			wantToken: "",
		},
		{
			name: "wrong password",
			request: &dto.LoginRequest{
				Login:    "testuser",
				Password: "wrongpassword",
			},
			setupMocks: func() (repository.UserRepositoryInterface, JWTServiceInterface) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				user := &models.User{
					ID:           "user-123",
					Login:        "testuser",
					PasswordHash: string(hashedPassword),
				}

				mockRepo := &mockUserRepository{
					users: map[string]*models.User{"testuser": user},
				}

				return mockRepo, &mockJWTService{}
			},
			wantErr:   errs.WrongLoginOrPassword,
			wantToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, mockJWT := tt.setupMocks()
			service := NewUserService(mockRepo, mockJWT)

			token, err := service.Login(context.Background(), tt.request)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}
