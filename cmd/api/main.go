package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Cxxlkid/go-task-api/internal/handler"
	"github.com/Cxxlkid/go-task-api/internal/repository"
	"github.com/Cxxlkid/go-task-api/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Chargement du .env
	godotenv.Load() // silencieux si pas de .env, normal en container

	// Logger structuré
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Connexion à la DB
	db, err := connectDB()
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("database connected")

	// Récupération des variables d'environnement
	jwtSecret := os.Getenv("JWT_SECRET")
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// Initialisation des couches (repository → usecase → handler)
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)
	taskUsecase := usecase.NewTaskUsecase(taskRepo, userRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	taskHandler := handler.NewTaskHandler(taskUsecase)

	// Router
	r := chi.NewRouter()

	// Middlewares globaux
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Routes publiques
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Post("/auth/register", userHandler.Register)
	r.Post("/auth/login", userHandler.Login)

	// Routes protégées par JWT
	r.Group(func(r chi.Router) {
		r.Use(handler.JWTMiddleware(jwtSecret))

		// Users
		r.Get("/users/me", userHandler.GetMe)

		// Tasks
		r.Post("/tasks", taskHandler.Create)
		r.Get("/tasks", taskHandler.List)
		r.Get("/tasks/{id}", taskHandler.GetByID)
		r.Put("/tasks/{id}", taskHandler.Update)
		r.Delete("/tasks/{id}", taskHandler.Delete)
	})

	// Démarrage du serveur
	slog.Info("server starting", "port", serverPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func connectDB() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
