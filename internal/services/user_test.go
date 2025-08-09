package services

import (
	"context"
	"testing"
	"time"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/storage"
	"github.com/Soliard/gophermart/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func Test_userService_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		req *dto.RegisterRequest
	}

	type setup struct {
		name string
		prep func(storage.UserRepositoryInterface) // подготовка данных
	}

	tests := []struct {
		name     string
		setup    setup
		args     args
		wantUser bool  // ожидаем ли пользователя
		wantErr  error // конкретная ошибка
	}{
		{
			name: "successful registration",
			setup: setup{
				name: "empty repository",
				prep: func(repo storage.UserRepositoryInterface) {
					// ничего не делаем - репозиторий пустой
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.RegisterRequest{
					Login:    "testuser",
					Password: "password123",
				},
			},
			wantUser: true,
			wantErr:  nil,
		},
		{
			name: "login already exists",
			setup: setup{
				name: "user already exists",
				prep: func(repo storage.UserRepositoryInterface) {
					// Создаем существующего пользователя
					existingUser := models.NewUser("testuser", "hashedpass")
					repo.Create(context.Background(), existingUser)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.RegisterRequest{
					Login:    "testuser",
					Password: "password123",
				},
			},
			wantUser: false,
			wantErr:  errs.LoginAlreadyExists,
		},
		{
			name: "empty login",
			setup: setup{
				name: "empty repository",
				prep: func(repo storage.UserRepositoryInterface) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.RegisterRequest{
					Login:    "",
					Password: "password123",
				},
			},
			wantUser: false,
			wantErr:  errs.EmptyLoginOrPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем чистые объекты для каждого теста
			userRepo := memory.NewUserRepository()
			jwtService := NewJWTService("test-secret", time.Hour)

			// Подготавливаем данные
			if tt.setup.prep != nil {
				tt.setup.prep(userRepo)
			}

			// Создаем сервис
			s := NewUserService(userRepo, jwtService)

			// Выполняем тест
			got, err := s.Register(tt.args.ctx, tt.args.req)

			// Проверяем ошибку
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)

			// Проверяем результат
			if tt.wantUser {
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.req.Login, got.Login)
				assert.NotEmpty(t, got.ID)
				assert.NotEmpty(t, got.PasswordHash)

				// Проверяем что пользователь действительно сохранился
				exists, _ := userRepo.UserExists(context.Background(), tt.args.req.Login)
				assert.True(t, exists)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
