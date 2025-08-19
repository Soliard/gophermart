package services

import (
	"context"
	"testing"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockWithdrawer struct {
	mock.Mock
}

type mockOrderService struct {
	mock.Mock
}

type mockBalanceService struct {
	mock.Mock
}

func (m *mockWithdrawer) Create(ctx context.Context, w *models.Withdrawal) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *mockWithdrawer) WithdrawalExists(ctx context.Context, userID string, orderNumber string) (bool, error) {
	args := m.Called(ctx, userID, orderNumber)
	return args.Bool(0), args.Error(1)
}

func (m *mockWithdrawer) GetWithdrawals(ctx context.Context, userID string) ([]*models.Withdrawal, error) {
	args := m.Called(ctx, userID)
	if v := args.Get(0); v != nil {
		return v.([]*models.Withdrawal), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOrderService) UploadOrder(ctx context.Context, userID, orderNumber string) (*models.Order, error) {
	args := m.Called(ctx, userID, orderNumber)
	if v := args.Get(0); v != nil {
		return v.(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOrderService) ValidateOrderNumber(ctx context.Context, orderNumber string) bool {
	args := m.Called(ctx, orderNumber)
	return args.Bool(0)
}

func (m *mockOrderService) GetUserOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	args := m.Called(ctx, userID)
	if v := args.Get(0); v != nil {
		return v.([]*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockBalanceService) GetBalance(ctx context.Context, userID string) (*models.Balance, error) {
	args := m.Called(ctx, userID)
	if v := args.Get(0); v != nil {
		return v.(*models.Balance), args.Error(1)
	}
	return nil, args.Error(1)
}

func Test_withdrawalService_ProcessWithdraw(t *testing.T) {
	type testCase struct {
		name          string
		userID        string
		order         string
		sum           float64
		setupMocks    func(*mockWithdrawer, *mockOrderService, *mockBalanceService)
		expectedError error
	}

	tests := []testCase{
		{
			name:   "успешное списание",
			userID: "user123",
			order:  "79927398713",
			sum:    100.50,
			setupMocks: func(mw *mockWithdrawer, mos *mockOrderService, mbs *mockBalanceService) {
				mos.On("ValidateOrderNumber", mock.Anything, "79927398713").Return(true)
				mw.On("WithdrawalExists", mock.Anything, "user123", "79927398713").Return(false, nil)
				mbs.On("GetBalance", mock.Anything, "user123").Return(&models.Balance{Current: 500}, nil)
				mw.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "недостаточно баллов",
			userID: "user123",
			order:  "79927398713",
			sum:    100.50,
			setupMocks: func(mw *mockWithdrawer, mos *mockOrderService, mbs *mockBalanceService) {
				mos.On("ValidateOrderNumber", mock.Anything, "79927398713").Return(true)
				mw.On("WithdrawalExists", mock.Anything, "user123", "79927398713").Return(false, nil)
				mbs.On("GetBalance", mock.Anything, "user123").Return(&models.Balance{Current: 99}, nil)
			},
			expectedError: errs.ErrBalanceInsufficient,
		},
		{
			name:   "невалидный номер запроса",
			userID: "user123",
			order:  "123",
			sum:    100.50,
			setupMocks: func(mw *mockWithdrawer, mos *mockOrderService, mbs *mockBalanceService) {
				mos.On("ValidateOrderNumber", mock.Anything, "123").Return(false)
			},
			expectedError: errs.ErrOrderIsNotValid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := new(mockWithdrawer)
			mos := new(mockOrderService)
			mbs := new(mockBalanceService)

			tt.setupMocks(mw, mos, mbs)

			service := NewWithdrawalService(mw, mbs, mos)
			err := service.ProcessWithdraw(context.Background(), tt.userID, tt.order, tt.sum)

			if tt.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mw.AssertExpectations(t)
			mos.AssertExpectations(t)
			mbs.AssertExpectations(t)
		})
	}
}
