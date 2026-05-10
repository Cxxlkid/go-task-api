package mocks

import (
	"github.com/Cxxlkid/go-task-api/internal/domain"
)

type UserRepositoryMock struct {
	CreateFn     func(user *domain.User) error
	GetByIDFn    func(id string) (*domain.User, error)
	GetByEmailFn func(email string) (*domain.User, error)
}

func (m *UserRepositoryMock) Create(user *domain.User) error {
	return m.CreateFn(user)
}

func (m *UserRepositoryMock) GetByID(id string) (*domain.User, error) {
	return m.GetByIDFn(id)
}

func (m *UserRepositoryMock) GetByEmail(email string) (*domain.User, error) {
	return m.GetByEmailFn(email)
}