package domain

import "time"

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

// Task is the core domain model
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

// CreateTaskRequest holds the payload for task creation
type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
}

// UpdateTaskRequest holds the payload for task update — all fields are optional
type UpdateTaskRequest struct {
	Title       *string     `json:"title,omitempty"`
	Description *string     `json:"description,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
	DueDate     *time.Time  `json:"due_date,omitempty"`
	AssignedTo  *string     `json:"assigned_to,omitempty"`
}

// TaskFilter holds the available filters and pagination options for listing tasks
type TaskFilter struct {
	Status     *TaskStatus `json:"status,omitempty"`
	AssignedTo *string     `json:"assigned_to,omitempty"`
	Page       int
	PageSize   int
	SortBy     string // "created_at" or "due_date"
}

// TaskRepository defines the persistence contract for tasks
type TaskRepository interface {
	Create(task *Task) error
	GetByID(id string) (*Task, error)
	Update(task *Task) error
	Delete(id string) error
	List(filter TaskFilter) ([]*Task, int, error)
}

// TaskUsecase defines the business logic contract for tasks
type TaskUsecase interface {
	Create(userID string, req CreateTaskRequest) (*Task, error)
	GetByID(id string) (*Task, error)
	Update(id string, req UpdateTaskRequest) (*Task, error)
	Delete(id string) error
	List(filter TaskFilter) ([]*Task, int, error)
}