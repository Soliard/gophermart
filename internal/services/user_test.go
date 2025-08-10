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

func setupUserService() (UserServiceInterface, storage.UserRepositoryInterface) {
	repo := memory.NewUserRepository()
	jwt := NewJWTService("secret", time.Hour)
	service := NewUserService(repo, jwt)
	return service, repo
}

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
			s, userRepo := setupUserService()

			// Подготавливаем данные
			if tt.setup.prep != nil {
				tt.setup.prep(userRepo)
			}

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
func Test_userService_Login(t *testing.T) {
	service, userRepo := setupUserService()
	type args struct {
		ctx context.Context
		req *dto.LoginRequest
	}

	type setup struct {
		name string
		prep func(storage.UserRepositoryInterface) // подготовка данных
	}

	tests := []struct {
		name      string
		setup     setup
		args      args
		wantToken bool  // ожидаем ли токен
		wantErr   error // конкретная ошибка
	}{
		{
			name: "successful login",
			setup: setup{
				name: "user exists with correct password",
				prep: func(repo storage.UserRepositoryInterface) {
					// Создаем пользователя с известным паролем
					// Используем Register для корректного хэширования
					service.Register(context.Background(), &dto.RegisterRequest{
						Login:    "testuser",
						Password: "password123",
					})
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Login:    "testuser",
					Password: "password123",
				},
			},
			wantToken: true,
			wantErr:   nil,
		},
		{
			name: "user not found",
			setup: setup{
				name: "empty repository",
				prep: func(repo storage.UserRepositoryInterface) {
					// пустой репозиторий
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Login:    "nonexistent",
					Password: "password123",
				},
			},
			wantToken: false,
			wantErr:   errs.WrongLoginOrPassword,
		},
		{
			name: "wrong password",
			setup: setup{
				name: "user exists with different password",
				prep: func(repo storage.UserRepositoryInterface) {
					// Создаем пользователя с одним паролем
					service, _ := setupUserService()
					service.Register(context.Background(), &dto.RegisterRequest{
						Login:    "testuser",
						Password: "correctpassword",
					})
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Login:    "testuser",
					Password: "wrongpassword",
				},
			},
			wantToken: false,
			wantErr:   errs.WrongLoginOrPassword,
		},
		{
			name: "empty login",
			setup: setup{
				name: "empty repository",
				prep: func(repo storage.UserRepositoryInterface) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Login:    "",
					Password: "password123",
				},
			},
			wantToken: false,
			wantErr:   errs.WrongLoginOrPassword, // пустой логин не найдется
		},
		{
			name: "empty password",
			setup: setup{
				name: "user exists",
				prep: func(repo storage.UserRepositoryInterface) {
					service, _ := setupUserService()
					service.Register(context.Background(), &dto.RegisterRequest{
						Login:    "testuser",
						Password: "password123",
					})
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Login:    "testuser",
					Password: "",
				},
			},
			wantToken: false,
			wantErr:   errs.WrongLoginOrPassword, // пустой пароль не совпадет с хэшем
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Подготавливаем данные
			if tt.setup.prep != nil {
				tt.setup.prep(userRepo)
			}

			// Выполняем тест
			token, err := service.Login(tt.args.ctx, tt.args.req)

			// Проверяем ошибку
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)

			// Проверяем результат
			if tt.wantToken {
				assert.NotEmpty(t, token)
				// проверить валидность токена?
			} else {
				assert.Empty(t, token)
			}
		})
	}
}
func Test_userService_RegisterAndLogin_Integration(t *testing.T) {
	s, _ := setupUserService()

	// 1. Регистрируем пользователя
	user, err := s.Register(context.Background(), &dto.RegisterRequest{
		Login:    "integrationuser",
		Password: "testpassword123",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// 2. Логинимся с теми же credentials
	token, err := s.Login(context.Background(), &dto.LoginRequest{
		Login:    "integrationuser",
		Password: "testpassword123",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 3. Пытаемся залогиниться с неправильным паролем
	_, err = s.Login(context.Background(), &dto.LoginRequest{
		Login:    "integrationuser",
		Password: "wrongpassword",
	})
	assert.Error(t, err)
	assert.Equal(t, errs.WrongLoginOrPassword, err)
}
