package memory

import "github.com/Soliard/gophermart/internal/repository"

type MemoryStorage struct {
	userRepository  repository.UserRepositoryInterface
	orderRepository repository.OrderRepositoryInterface
}

func (s *MemoryStorage) UserRepository() repository.UserRepositoryInterface {
	return s.userRepository
}

func (s *MemoryStorage) OrderRepository() repository.OrderRepositoryInterface {
	return s.orderRepository
}

func NewMemoryStorage() (repository.Storage, error) {
	return nil, nil
	// return &MemoryStorage{
	// 	userRepository:  newUserRepository(),
	// 	orderRepository: newOrderRepository(),
	// }, nil
}
