package handlers

import (
	"encoding/json"
	"net/http"

	"daycare/internal/service"
)

type AdminUsersHandler struct {
	svc *service.AdminUsersService
}

func NewAdminUsersHandler(svc *service.AdminUsersService) *AdminUsersHandler {
	return &AdminUsersHandler{svc: svc}
}

type createUserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // ADMIN/STAFF
}

func (h *AdminUsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	id, err := h.svc.CreateUser(r.Context(), service.CreateUserInput{
		Name: req.Name, Email: req.Email, Password: req.Password, Role: req.Role,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid input"})
		case service.ErrConflict:
			writeJSON(w, http.StatusConflict, map[string]any{"error": "email already exists"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"user_id": id})
}
