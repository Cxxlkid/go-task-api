package domain

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Ce qu'on reçoit pour créer une tâche
type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
}

// Ce qu'on reçoit pour modifier une tâche
type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
}

// Filtres pour lister les tâches
type TaskFilter struct {
	Status     *TaskStatus `json:"status,omitempty"`
	AssignedTo *string     `json:"assigned_to,omitempty"`
	Page       int
	PageSize   int
	SortBy     string // "created_at" ou "due_date"
}

// L'interface que le repository DOIT implémenter
type TaskRepository interface {
	Create(task *Task) error
	GetByID(id string) (*Task, error)
	Update(task *Task) error
	Delete(id string) error
	List(filter TaskFilter) ([]*Task, int, error)
}

// L'interface que le usecase DOIT implémenter
type TaskUsecase interface {
	Create(userID string, req CreateTaskRequest) (*Task, error)
	GetByID(id string) (*Task, error)
	Update(id string, req UpdateTaskRequest) (*Task, error)
	Delete(id string) error
	List(filter TaskFilter) ([]*Task, int, error)
}