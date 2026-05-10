package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Cxxlkid/go-task-api/internal/domain"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userUsecase.Register(req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	authResponse, err := h.userUsecase.Login(req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, authResponse)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	user, err := h.userUsecase.GetByID(userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}