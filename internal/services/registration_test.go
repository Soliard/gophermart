package services

import (
	"context"
	"testing"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockUserRegistrator struct {
	mock.Mock
}

func (m *mockUserRegistrator) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRegistrator) UserExists(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

func Test_registrationService_Register(t *testing.T) {
	tests := []struct {
		name          string
		regReq        dto.RegisterRequest
		mockSetup     func(*mockUserRegistrator)
		expectedError error
	}{
		{
			name:          "новый юзер",
			regReq:        dto.RegisterRequest{Login: "u1", Password: "123"},
			expectedError: nil,
			mockSetup: func(m *mockUserRegistrator) {
				m.On("UserExists", mock.Anything, "u1").Return(false, nil)
				m.On("Create", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
					return u.Login == "u1"
				})).Return(nil)
			},
		},
		{
			name:          "юзер уже существует",
			regReq:        dto.RegisterRequest{Login: "u2", Password: "asdf"},
			expectedError: errs.ErrUserAlreadyExists,
			mockSetup: func(mur *mockUserRegistrator) {
				mur.On("UserExists", mock.Anything, "u2").Return(true, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(mockUserRegistrator)
			tt.mockSetup(mock)
			ctx := context.Background()
			service := NewRegistrationService(mock)
			user, err := service.Register(ctx, &tt.regReq)

			if err == nil {
				require.NotNil(t, user)
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.expectedError, err)
				require.Nil(t, user)
			}

			mock.AssertExpectations(t)
		})
	}
}
