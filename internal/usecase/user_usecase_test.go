package usecase_test

import (
	"errors"
	"testing"

	"github.com/Cxxlkid/go-task-api/internal/domain"
	"github.com/Cxxlkid/go-task-api/internal/mocks"
	"github.com/Cxxlkid/go-task-api/internal/usecase"
)

func TestRegister_Success(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetByEmailFn: func(email string) (*domain.User, error) {
			return nil, errors.New("not found") // pas d'user existant
		},
		CreateFn: func(user *domain.User) error {
			return nil
		},
	}

	uc := usecase.NewUserUsecase(repo, "secret")
	user, err := uc.Register(domain.CreateUserRequest{
		Email:    "test@test.com",
		Password: "password123",
		Name:     "Test User",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "test@test.com" {
		t.Errorf("expected email test@test.com, got %s", user.Email)
	}
	if user.PasswordHash == "" {
		t.Error("expected password hash to be set")
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetByEmailFn: func(email string) (*domain.User, error) {
			return &domain.User{Email: email}, nil // user déjà existant
		},
	}

	uc := usecase.NewUserUsecase(repo, "secret")
	_, err := uc.Register(domain.CreateUserRequest{
		Email:    "test@test.com",
		Password: "password123",
		Name:     "Test User",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "email already in use" {
		t.Errorf("expected 'email already in use', got %s", err.Error())
	}
}

func TestRegister_InvalidInput(t *testing.T) {
	repo := &mocks.UserRepositoryMock{}
	uc := usecase.NewUserUsecase(repo, "secret")

	_, err := uc.Register(domain.CreateUserRequest{
		Email:    "not-an-email",
		Password: "123",
		Name:     "T",
	})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	// On crée d'abord un vrai hash pour le test
	registerRepo := &mocks.UserRepositoryMock{
		GetByEmailFn: func(email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(user *domain.User) error { return nil },
	}
	uc := usecase.NewUserUsecase(registerRepo, "secret")
	registered, _ := uc.Register(domain.CreateUserRequest{
		Email:    "test@test.com",
		Password: "password123",
		Name:     "Test User",
	})

	// Maintenant on teste le login avec le vrai hash
	loginRepo := &mocks.UserRepositoryMock{
		GetByEmailFn: func(email string) (*domain.User, error) {
			return registered, nil
		},
	}
	uc2 := usecase.NewUserUsecase(loginRepo, "secret")
	resp, err := uc2.Login(domain.LoginRequest{
		Email:    "test@test.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token to be set")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	registerRepo := &mocks.UserRepositoryMock{
		GetByEmailFn: func(email string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(user *domain.User) error { return nil },
	}
	uc := usecase.NewUserUsecase(registerRepo, "secret")
	registered, _ := uc.Register(domain.CreateUserRequest{
		Email:    "test@test.com",
		Password: "password123",
		Name:     "Test User",
	})

	loginRepo := &mocks.UserRepositoryMock{
		GetByEmailFn: func(email string) (*domain.User, error) {
			return registered, nil
		},
	}
	uc2 := usecase.NewUserUsecase(loginRepo, "secret")
	_, err := uc2.Login(domain.LoginRequest{
		Email:    "test@test.com",
		Password: "wrongpassword",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}