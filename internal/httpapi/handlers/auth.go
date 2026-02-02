package handlers

import (
	"encoding/json"
	"net/http"

	"daycare/internal/domain"
	"daycare/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler { return &AuthHandler{auth: a} }

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	res, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid credentials"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": res.Token,
		"user": map[string]any{
			"id":    res.User.ID,
			"name":  res.User.Name,
			"email": res.User.Email,
			"role":  res.User.Role,
		},
	})
	_ = domain.ErrUnauthorized
}
