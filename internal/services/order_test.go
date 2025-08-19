package services

import (
	"context"
	"testing"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockOrderCreator struct {
	mock.Mock
}

func (m *mockOrderCreator) Create(ctx context.Context, o *models.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOrderCreator) GetByNumber(ctx context.Context, number string) (*models.Order, error) {
	args := m.Called(ctx, number)
	if v := args.Get(0); v != nil {
		return v.(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOrderCreator) GetUserOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	args := m.Called(ctx, userID)
	if v := args.Get(0); v != nil {
		return v.([]*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestOrderService_UploadOrder(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		orderNumber   string
		mockSetup     func(*mockOrderCreator)
		expectedError error
	}{
		{
			name:        "успешная загрузка нового заказа",
			userID:      "user123",
			orderNumber: "79927398713",
			mockSetup: func(m *mockOrderCreator) {
				m.On("GetByNumber", mock.Anything, "79927398713").
					Return(nil, errs.ErrOrderNotFound)

				m.On("Create", mock.Anything, mock.MatchedBy(func(o *models.Order) bool {
					return o.Number == "79927398713" && o.UserID == "user123"
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "заказ уже загружен другим пользователем",
			userID:      "user123",
			orderNumber: "79927398713",
			mockSetup: func(m *mockOrderCreator) {
				existingOrder := &models.Order{
					Number: "79927398713",
					UserID: "user456",
				}
				m.On("GetByNumber", mock.Anything, "79927398713").
					Return(existingOrder, nil)
			},
			expectedError: errs.ErrOrderAlreadyUploadedByOtherUser,
		},
		{
			name:        "заказ уже загружен этим пользователем",
			userID:      "user123",
			orderNumber: "79927398713",
			mockSetup: func(m *mockOrderCreator) {
				existingOrder := &models.Order{
					Number: "79927398713",
					UserID: "user123",
				}
				m.On("GetByNumber", mock.Anything, "79927398713").
					Return(existingOrder, nil)
			},
			expectedError: errs.ErrOrderAlreadyUploadedByThisUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockOrderCreator)

			tt.mockSetup(mockRepo)

			service := NewOrderService(mockRepo)

			ctx := context.Background()
			result, err := service.UploadOrder(ctx, tt.userID, tt.orderNumber)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_ValidateOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderNumber string
		isValid     bool
	}{
		{
			name:        "валидный номер заказа",
			orderNumber: "79927398713",
			isValid:     true,
		},
		{
			name:        "невалидный номер заказа",
			orderNumber: "123",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockOrderCreator)
			service := NewOrderService(mockRepo)
			ctx := context.Background()
			require.Equal(t, tt.isValid, service.ValidateOrderNumber(ctx, tt.orderNumber))
		})
	}
}
