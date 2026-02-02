package handlers

import (
	"encoding/json"
	"net/http"

	"daycare/internal/service"
)

type BootstrapHandler struct {
	svc    *service.BootstrapService
	appEnv string
}

func NewBootstrapHandler(svc *service.BootstrapService, appEnv string) *BootstrapHandler {
	return &BootstrapHandler{svc: svc, appEnv: appEnv}
}

type bootstrapReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *BootstrapHandler) CreateFirstAdmin(w http.ResponseWriter, r *http.Request) {
	// Solo dev (para no dejar un hueco en prod)
	if h.appEnv != "dev" {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}

	var req bootstrapReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	id, err := h.svc.CreateFirstAdmin(r.Context(), service.BootstrapAdminInput{
		Name: req.Name, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid input"})
		case service.ErrConflict:
			writeJSON(w, http.StatusConflict, map[string]any{"error": "admin already exists"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"admin_user_id": id})
}
