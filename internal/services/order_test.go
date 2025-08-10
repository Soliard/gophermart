package services

import (
	"context"
	"testing"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/storage"
	"github.com/Soliard/gophermart/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func setupOrderService() (OrderServiceInterface, storage.OrderRepositoryInterface) {
	repo := memory.NewOrderRepository()
	service := NewOrderService(repo)
	return service, repo
}

func Test_orderService_UploadOrder(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		orderNumber string
		prepareRepo func(storage.OrderRepositoryInterface)
		wantErr     error
		checkResult func(*testing.T, *models.Order)
	}{
		{
			name:        "successful upload new order",
			userID:      "user123",
			orderNumber: "12345678903",
			prepareRepo: func(repo storage.OrderRepositoryInterface) {
				// пустой репозиторий
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
			prepareRepo: func(repo storage.OrderRepositoryInterface) {
				// Создаем заказ от того же пользователя
				existingOrder := models.NewOrder("12345678903", "user123")
				repo.Create(context.Background(), existingOrder)
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
			prepareRepo: func(repo storage.OrderRepositoryInterface) {
				// Создаем заказ от другого пользователя
				existingOrder := models.NewOrder("12345678903", "user456")
				repo.Create(context.Background(), existingOrder)
			},
			wantErr: errs.OrderAlreadyUploadedByOtherUser,
			checkResult: func(t *testing.T, order *models.Order) {
				assert.Nil(t, order)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo := setupOrderService()

			// Подготавливаем данные
			if tt.prepareRepo != nil {
				tt.prepareRepo(repo)
			}

			// Выполняем тест
			got, err := service.UploadOrder(context.Background(), tt.userID, tt.orderNumber)

			// Проверяем ошибку
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}

			// Проверяем результат
			if tt.checkResult != nil {
				tt.checkResult(t, got)
			}

			// Дополнительная проверка: если заказ должен был сохраниться
			if tt.wantErr == nil {
				savedOrder, err := repo.GetByNumber(context.Background(), tt.orderNumber)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, savedOrder.UserID)
			}
		})
	}
}

func Test_orderService_ValidateOrderNumber(t *testing.T) {
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
			service, _ := setupOrderService()

			got := service.ValidateOrderNumber(context.Background(), tt.orderNumber)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_orderService_Integration(t *testing.T) {
	service, _ := setupOrderService()

	// Тест полного флоу
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
