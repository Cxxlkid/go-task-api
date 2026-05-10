package usecase_test

import (
	"errors"
	"testing"

	"github.com/Cxxlkid/go-task-api/internal/domain"
	"github.com/Cxxlkid/go-task-api/internal/mocks"
	"github.com/Cxxlkid/go-task-api/internal/usecase"
)

func TestCreateTask_Success(t *testing.T) {
	taskRepo := &mocks.TaskRepositoryMock{
		CreateFn: func(task *domain.Task) error { return nil },
	}
	userRepo := &mocks.UserRepositoryMock{}

	uc := usecase.NewTaskUsecase(taskRepo, userRepo)
	task, err := uc.Create("user-123", domain.CreateTaskRequest{
		Title:       "Ma tâche",
		Description: "Description",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.Title != "Ma tâche" {
		t.Errorf("expected title 'Ma tâche', got %s", task.Title)
	}
	if task.Status != domain.StatusTodo {
		t.Errorf("expected status todo, got %s", task.Status)
	}
	if task.CreatedBy != "user-123" {
		t.Errorf("expected created_by user-123, got %s", task.CreatedBy)
	}
}

func TestCreateTask_InvalidTitle(t *testing.T) {
	taskRepo := &mocks.TaskRepositoryMock{}
	userRepo := &mocks.UserRepositoryMock{}

	uc := usecase.NewTaskUsecase(taskRepo, userRepo)
	_, err := uc.Create("user-123", domain.CreateTaskRequest{
		Title: "ab", // trop court
	})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestCreateTask_AssignedUserNotFound(t *testing.T) {
	assignedID := "unknown-user"
	taskRepo := &mocks.TaskRepositoryMock{}
	userRepo := &mocks.UserRepositoryMock{
		GetByIDFn: func(id string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
	}

	uc := usecase.NewTaskUsecase(taskRepo, userRepo)
	_, err := uc.Create("user-123", domain.CreateTaskRequest{
		Title:      "Ma tâche",
		AssignedTo: &assignedID,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "assigned user not found" {
		t.Errorf("expected 'assigned user not found', got %s", err.Error())
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	taskRepo := &mocks.TaskRepositoryMock{
		GetByIDFn: func(id string) (*domain.Task, error) {
			return nil, errors.New("not found")
		},
	}
	userRepo := &mocks.UserRepositoryMock{}

	uc := usecase.NewTaskUsecase(taskRepo, userRepo)
	err := uc.Delete("unknown-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}