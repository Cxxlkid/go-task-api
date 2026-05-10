package usecase

import (
	"fmt"
	"time"

	"github.com/Cxxlkid/go-task-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/Cxxlkid/go-task-api/internal/validator"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewUserUsecase(userRepo domain.UserRepository, jwtSecret string) domain.UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) Register(req domain.CreateUserRequest) (*domain.User, error) {
	// Validation
	if errs := validator.ValidateRegister(req.Email, req.Password, req.Name); len(errs) > 0 {
		return nil, errs
	}
	// Vérifier si l'email existe déjà
	existing, _ := u.userRepo.GetByEmail(req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already in use")
	}

	// Hasher le mot de passe
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("userUsecase.Register hash: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("userUsecase.Register: %w", err)
	}

	return user, nil
}

func (u *userUsecase) Login(req domain.LoginRequest) (*domain.AuthResponse, error) {
	// Validation
	if errs := validator.ValidateLogin(req.Email, req.Password); len(errs) > 0 {
		return nil, errs
	}
	// Récupérer le user par email
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Vérifier le mot de passe
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Générer le JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("userUsecase.Login token: %w", err)
	}

	return &domain.AuthResponse{
		Token: tokenString,
		User:  *user,
	}, nil
}

func (u *userUsecase) GetByID(id string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("userUsecase.GetByID: %w", err)
	}
	return user, nil
}