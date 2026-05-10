package domain

import "time"

// User is the core domain model
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never exposed in JSON responses
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateUserRequest holds the payload for user registration
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest holds the payload for user authentication
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse holds the JWT token and user info returned after a successful login
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UserRepository defines the persistence contract for users
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
}

// UserUsecase defines the business logic contract for users
type UserUsecase interface {
	Register(req CreateUserRequest) (*User, error)
	Login(req LoginRequest) (*AuthResponse, error)
	GetByID(id string) (*User, error)
}