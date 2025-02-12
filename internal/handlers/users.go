// internal/handlers/users.go
package handlers

import (
    "encoding/json"
    "net/http"
    "print-automation/internal/repository"
)

type UserHandler struct {
    repo *repository.UserRepository
}

type CreateUserRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
    return &UserHandler{repo: repo}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // В реальном приложении здесь должно быть хеширование пароля
    user, err := h.repo.Create(req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}