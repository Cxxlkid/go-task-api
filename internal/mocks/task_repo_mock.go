package mocks

import (
	"github.com/Cxxlkid/go-task-api/internal/domain"
)

type TaskRepositoryMock struct {
	CreateFn  func(task *domain.Task) error
	GetByIDFn func(id string) (*domain.Task, error)
	UpdateFn  func(task *domain.Task) error
	DeleteFn  func(id string) error
	ListFn    func(filter domain.TaskFilter) ([]*domain.Task, int, error)
}

func (m *TaskRepositoryMock) Create(task *domain.Task) error {
	return m.CreateFn(task)
}

func (m *TaskRepositoryMock) GetByID(id string) (*domain.Task, error) {
	return m.GetByIDFn(id)
}

func (m *TaskRepositoryMock) Update(task *domain.Task) error {
	return m.UpdateFn(task)
}

func (m *TaskRepositoryMock) Delete(id string) error {
	return m.DeleteFn(id)
}

func (m *TaskRepositoryMock) List(filter domain.TaskFilter) ([]*domain.Task, int, error) {
	return m.ListFn(filter)
}