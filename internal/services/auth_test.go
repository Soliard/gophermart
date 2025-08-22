package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserLoginner struct {
	mock.Mock
}

type mockJWTService struct {
	mock.Mock
}

func (m *mockUserLoginner) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	args := m.Called(ctx, login)
	if v := args.Get(0); v != nil {
		return v.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserLoginner) UpdateLoginTime(ctx context.Context, userID string, time time.Time) error {
	args := m.Called(ctx, userID, time)
	return args.Error(0)
}

func (m *mockJWTService) GenerateToken(u *models.User) (string, error) {
	args := m.Called(u)
	return args.String(0), args.Error(1)
}

func (m *mockJWTService) GetClaims(tokenString string) (*UserContext, error) {
	args := m.Called(tokenString)
	if v := args.Get(0); v != nil {
		return v.(*UserContext), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	validPassword := "correct_password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)

	testUser := &models.User{
		ID:           "user123",
		Login:        "testuser",
		PasswordHash: string(hashedPassword),
		Roles:        []models.Role{models.RoleUser},
	}

	tests := []struct {
		name          string
		loginReq      *dto.LoginRequest
		mockSetup     func(*mockUserLoginner, *mockJWTService)
		expectedToken string
		expectedError error
	}{
		{
			name: "успешный логин",
			loginReq: &dto.LoginRequest{
				Login:    "testuser",
				Password: validPassword,
			},
			mockSetup: func(mul *mockUserLoginner, mjwt *mockJWTService) {
				mul.On("GetByLogin", mock.Anything, "testuser").
					Return(testUser, nil)

				mjwt.On("GenerateToken", testUser).
					Return("valid-jwt-token", nil)

				mul.On("UpdateLoginTime", mock.Anything, "user123", mock.MatchedBy(func(t time.Time) bool {
					return !t.IsZero()
				})).Return(nil)
			},
			expectedToken: "valid-jwt-token",
			expectedError: nil,
		},
		{
			name: "неверный пароль",
			loginReq: &dto.LoginRequest{
				Login:    "testuser",
				Password: "wrong_password",
			},
			mockSetup: func(mul *mockUserLoginner, mjwt *mockJWTService) {
				mul.On("GetByLogin", mock.Anything, "testuser").
					Return(testUser, nil)
			},
			expectedToken: "",
			expectedError: errs.ErrWrongLoginOrPassword,
		},
		{
			name: "пользователь не найден",
			loginReq: &dto.LoginRequest{
				Login:    "nonexistent",
				Password: "anypassword",
			},
			mockSetup: func(mul *mockUserLoginner, mjwt *mockJWTService) {
				mul.On("GetByLogin", mock.Anything, "nonexistent").
					Return(nil, errs.ErrUserNotFound)
			},
			expectedToken: "",
			expectedError: errs.ErrWrongLoginOrPassword,
		},
		{
			name: "ошибка базы данных при поиске пользователя",
			loginReq: &dto.LoginRequest{
				Login:    "testuser",
				Password: validPassword,
			},
			mockSetup: func(mul *mockUserLoginner, mjwt *mockJWTService) {
				mul.On("GetByLogin", mock.Anything, "testuser").
					Return(nil, errors.New("database error"))
			},
			expectedToken: "",
			expectedError: errors.New("database error"),
		},
		{
			name: "ошибка генерации токена",
			loginReq: &dto.LoginRequest{
				Login:    "testuser",
				Password: validPassword,
			},
			mockSetup: func(mul *mockUserLoginner, mjwt *mockJWTService) {
				mul.On("GetByLogin", mock.Anything, "testuser").
					Return(testUser, nil)

				mjwt.On("GenerateToken", testUser).
					Return("", errors.New("jwt generation error"))

			},
			expectedToken: "",
			expectedError: errors.New("jwt generation error"),
		},
		{
			name: "ошибка обновления времени логина",
			loginReq: &dto.LoginRequest{
				Login:    "testuser",
				Password: validPassword,
			},
			mockSetup: func(mul *mockUserLoginner, mjwt *mockJWTService) {
				mul.On("GetByLogin", mock.Anything, "testuser").
					Return(testUser, nil)

				mjwt.On("GenerateToken", testUser).
					Return("valid-token", nil)

				mul.On("UpdateLoginTime", mock.Anything, "user123", mock.Anything).
					Return(errors.New("update failed"))
			},
			expectedToken: "valid-token",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(mockUserLoginner)
			mockJWT := new(mockJWTService)

			tt.mockSetup(mockUserRepo, mockJWT)

			service := NewAuthService(mockUserRepo, mockJWT)

			ctx := context.Background()
			token, err := service.Login(ctx, tt.loginReq)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedToken, token)

			mockUserRepo.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
		})
	}
}
