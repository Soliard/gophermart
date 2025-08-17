// internal/services/order_test.go
package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/repository"
	"github.com/stretchr/testify/assert"
)

type mockOrderRepository struct {
	orders           map[string]*models.Order
	shouldFailCreate bool
	shouldFailGet    bool
	getResult        *models.Order
	getError         error
}

func (m *mockOrderRepository) Create(ctx context.Context, order *models.Order) error {
	if m.shouldFailCreate {
		return errors.New("create error")
	}

	if m.orders == nil {
		m.orders = make(map[string]*models.Order)
	}
	m.orders[order.Number] = order
	return nil
}

func (m *mockOrderRepository) GetByNumber(ctx context.Context, number string) (*models.Order, error) {
	if m.shouldFailGet {
		return nil, errors.New("get error")
	}

	if m.getError != nil {
		return nil, m.getError
	}

	if m.getResult != nil {
		return m.getResult, nil
	}

	order, exists := m.orders[number]
	if !exists {
		return nil, errs.OrderNotFound
	}
	return order, nil
}

func TestOrderService_UploadOrder(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		orderNumber string
		setupMock   func() repository.OrderRepositoryInterface
		wantErr     error
		checkResult func(*testing.T, *models.Order)
	}{
		{
			name:        "successful upload new order",
			userID:      "user123",
			orderNumber: "12345678903",
			setupMock: func() repository.OrderRepositoryInterface {
				return &mockOrderRepository{
					getError: errs.OrderNotFound, // заказ не найден - можно создавать
				}
			},
			wantErr: nil,
			checkResult: func(t *testing.T, order *models.Order) {
				assert.NotNil(t, order)
				assert.Equal(t, "12345678903", order.Number)
				assert.Equal(t, "user123", order.UserID)
				assert.Equal(t, models.StatusNew, order.Status)
				assert.Nil(t, order.Accrual)
				assert.False(t, order.UploadedAt.IsZero())
			},
		},
		{
			name:        "order already uploaded by same user",
			userID:      "user123",
			orderNumber: "12345678903",
			setupMock: func() repository.OrderRepositoryInterface {
				existingOrder := &models.Order{
					Number: "12345678903",
					UserID: "user123",
					Status: models.StatusNew,
				}
				return &mockOrderRepository{
					getResult: existingOrder, // заказ уже существует от того же пользователя
				}
			},
			wantErr: errs.OrderAlreadyUploadedByThisUser,
			checkResult: func(t *testing.T, order *models.Order) {
				assert.Nil(t, order)
			},
		},
		{
			name:        "order already uploaded by other user",
			userID:      "user123",
			orderNumber: "12345678903",
			setupMock: func() repository.OrderRepositoryInterface {
				existingOrder := &models.Order{
					Number: "12345678903",
					UserID: "user456", // другой пользователь
					Status: models.StatusNew,
				}
				return &mockOrderRepository{
					getResult: existingOrder,
				}
			},
			wantErr: errs.OrderAlreadyUploadedByOtherUser,
			checkResult: func(t *testing.T, order *models.Order) {
				assert.Nil(t, order)
			},
		},
		{
			name:        "repository error on GetByNumber",
			userID:      "user123",
			orderNumber: "12345678903",
			setupMock: func() repository.OrderRepositoryInterface {
				return &mockOrderRepository{
					shouldFailGet: true,
				}
			},
			wantErr: errors.New("get error"),
			checkResult: func(t *testing.T, order *models.Order) {
				assert.Nil(t, order)
			},
		},
		{
			name:        "repository error on Create",
			userID:      "user123",
			orderNumber: "12345678903",
			setupMock: func() repository.OrderRepositoryInterface {
				return &mockOrderRepository{
					getError:         errs.OrderNotFound, // заказ не найден
					shouldFailCreate: true,               // но создание падает
				}
			},
			wantErr: errors.New("create error"),
			checkResult: func(t *testing.T, order *models.Order) {
				assert.Nil(t, order)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок
			mockRepo := tt.setupMock()

			// Создаем сервис
			service := NewOrderService(mockRepo)

			// Выполняем тест
			order, err := service.UploadOrder(context.Background(), tt.userID, tt.orderNumber)

			// Проверяем ошибку
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем результат
			tt.checkResult(t, order)
		})
	}
}

func TestOrderService_ValidateOrderNumber(t *testing.T) {
	tests := []struct {
		name        string
		orderNumber string
		want        bool
	}{
		{
			name:        "valid order number",
			orderNumber: "12345678903", // валидный по алгоритму Луна
			want:        true,
		},
		{
			name:        "invalid order number - wrong checksum",
			orderNumber: "12345678904", // невалидный по алгоритму Луна
			want:        false,
		},
		{
			name:        "invalid order number - not numeric",
			orderNumber: "abc123",
			want:        false,
		},
		{
			name:        "empty order number",
			orderNumber: "",
			want:        false,
		},
		{
			name:        "single digit",
			orderNumber: "0",
			want:        false,
		},
		{
			name:        "another valid number",
			orderNumber: "4532015112830366", // валидный номер карты
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Для ValidateOrderNumber не нужен репозиторий
			mockRepo := &mockOrderRepository{}
			service := NewOrderService(mockRepo)

			got := service.ValidateOrderNumber(context.Background(), tt.orderNumber)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Интеграционный тест для проверки полного флоу
func TestOrderService_Integration(t *testing.T) {
	mockRepo := &mockOrderRepository{}
	service := NewOrderService(mockRepo)

	orderNumber := "12345678903"
	userID := "user123"

	// 1. Валидируем номер
	isValid := service.ValidateOrderNumber(context.Background(), orderNumber)
	assert.True(t, isValid)

	// 2. Загружаем заказ
	order, err := service.UploadOrder(context.Background(), userID, orderNumber)
	assert.NoError(t, err)
	assert.NotNil(t, order)

	// 3. Пытаемся загрузить тот же заказ тем же пользователем
	_, err = service.UploadOrder(context.Background(), userID, orderNumber)
	assert.Equal(t, errs.OrderAlreadyUploadedByThisUser, err)

	// 4. Пытаемся загрузить тот же заказ другим пользователем
	_, err = service.UploadOrder(context.Background(), "user456", orderNumber)
	assert.Equal(t, errs.OrderAlreadyUploadedByOtherUser, err)
}
