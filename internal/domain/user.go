package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // le "-" = jamais renvoyé en JSON
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Ce qu'on reçoit pour créer un user
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Ce qu'on reçoit pour se connecter
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Ce qu'on renvoie après login
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// L'interface que le repository DOIT implémenter
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
}

// L'interface que le usecase DOIT implémenter
type UserUsecase interface {
	Register(req CreateUserRequest) (*User, error)
	Login(req LoginRequest) (*AuthResponse, error)
	GetByID(id string) (*User, error)
}