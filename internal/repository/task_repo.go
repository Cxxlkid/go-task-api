package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cxxlkid/go-task-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(context.Background(), query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.AssignedTo,
		task.CreatedBy,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("taskRepository.Create: %w", err)
	}
	return nil
}

func (r *taskRepository) GetByID(id string) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at
		FROM tasks WHERE id = $1
	`
	task := &domain.Task{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DueDate,
		&task.AssignedTo,
		&task.CreatedBy,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.GetByID: %w", err)
	}
	return task, nil
}

func (r *taskRepository) Update(task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, due_date = $4, assigned_to = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(context.Background(), query,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.AssignedTo,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("taskRepository.Update: %w", err)
	}
	return nil
}

func (r *taskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("taskRepository.Delete: %w", err)
	}
	return nil
}

func (r *taskRepository) List(filter domain.TaskFilter) ([]*domain.Task, int, error) {
	// Construction dynamique de la requête selon les filtres
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.AssignedTo != nil {
		conditions = append(conditions, fmt.Sprintf("assigned_to = $%d", argIndex))
		args = append(args, *filter.AssignedTo)
		argIndex++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Tri
	sortBy := "created_at"
	if filter.SortBy == "due_date" {
		sortBy = "due_date"
	}

	// Pagination
	if filter.PageSize == 0 {
		filter.PageSize = 10
	}
	if filter.Page == 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize

	// Compte total pour la pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks %s", where)
	var total int
	err := r.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("taskRepository.List count: %w", err)
	}

	// Requête principale
	query := fmt.Sprintf(`
		SELECT id, title, description, status, due_date, assigned_to, created_by, created_at, updated_at
		FROM tasks %s
		ORDER BY %s DESC
		LIMIT $%d OFFSET $%d
	`, where, sortBy, argIndex, argIndex+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("taskRepository.List: %w", err)
	}
	defer rows.Close()

	tasks := []*domain.Task{}
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.DueDate,
			&task.AssignedTo,
			&task.CreatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("taskRepository.List scan: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}