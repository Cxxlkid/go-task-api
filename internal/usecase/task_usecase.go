package usecase

import (
	"fmt"
	"time"

	"github.com/Cxxlkid/go-task-api/internal/domain"
	"github.com/google/uuid"
	"github.com/Cxxlkid/go-task-api/internal/validator"

)

type taskUsecase struct {
	taskRepo domain.TaskRepository
	userRepo domain.UserRepository
}

func NewTaskUsecase(taskRepo domain.TaskRepository, userRepo domain.UserRepository) domain.TaskUsecase {
	return &taskUsecase{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (u *taskUsecase) Create(userID string, req domain.CreateTaskRequest) (*domain.Task, error) {
	// Validation
	if errs := validator.ValidateCreateTask(req.Title); len(errs) > 0 {
		return nil, errs
	}
	// Vérifier que l'user assigné existe si fourni
	if req.AssignedTo != nil {
		_, err := u.userRepo.GetByID(*req.AssignedTo)
		if err != nil {
			return nil, fmt.Errorf("assigned user not found")
		}
	}

	now := time.Now()
	task := &domain.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.StatusTodo,
		DueDate:     req.DueDate,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("taskUsecase.Create: %w", err)
	}

	return task, nil
}

func (u *taskUsecase) GetByID(id string) (*domain.Task, error) {
	task, err := u.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("taskUsecase.GetByID: %w", err)
	}
	return task, nil
}

func (u *taskUsecase) Update(id string, req domain.UpdateTaskRequest) (*domain.Task, error) {
	// Validation
	status := (*string)(req.Status)
	if errs := validator.ValidateUpdateTask(status); len(errs) > 0 {
		return nil, errs
	}
	task, err := u.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	// On met à jour uniquement les champs fournis
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.AssignedTo != nil {
		_, err := u.userRepo.GetByID(*req.AssignedTo)
		if err != nil {
			return nil, fmt.Errorf("assigned user not found")
		}
		task.AssignedTo = req.AssignedTo
	}

	task.UpdatedAt = time.Now()

	if err := u.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("taskUsecase.Update: %w", err)
	}

	return task, nil
}

func (u *taskUsecase) Delete(id string) error {
	_, err := u.taskRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	return u.taskRepo.Delete(id)
}

func (u *taskUsecase) List(filter domain.TaskFilter) ([]*domain.Task, int, error) {
	tasks, total, err := u.taskRepo.List(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("taskUsecase.List: %w", err)
	}
	return tasks, total, nil
}