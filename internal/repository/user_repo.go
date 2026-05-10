package repository

import (
	"context"
	"fmt"

	"github.com/Cxxlkid/go-task-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(context.Background(), query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("userRepository.Create: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByID: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByEmail: %w", err)
	}
	return user, nil
}