package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"daycare/internal/domain"
	"daycare/internal/service"

	"github.com/go-chi/chi/v5"
)

type ChildrenHandler struct {
	svc *service.ChildrenService
}

func NewChildrenHandler(s *service.ChildrenService) *ChildrenHandler { return &ChildrenHandler{svc: s} }

type createChildReq struct {
	DocumentNumber string `json:"document_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	GuardianName   string `json:"guardian_name"`
	GuardianPhone  string `json:"guardian_phone"`
}

func (h *ChildrenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createChildReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	child, err := h.svc.Create(r.Context(), service.CreateChildInput(req))
	if err != nil {
		switch err {
		case domain.ErrInvalid:
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing required fields"})
		case domain.ErrConflict:
			writeJSON(w, http.StatusConflict, map[string]any{"error": "child already exists"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":              child.ID,
		"document_number": child.DocumentNumber,
		"first_name":      child.FirstName,
		"last_name":       child.LastName,
		"guardian_name":   child.GuardianName,
		"guardian_phone":  child.GuardianPhone,
		"created_at":      child.CreatedAt,
	})
}

func (h *ChildrenHandler) GetByDocument(w http.ResponseWriter, r *http.Request) {
	doc := chi.URLParam(r, "document")
	child, err := h.svc.GetByDocument(r.Context(), doc)
	if err != nil {
		if err == domain.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "child not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		return
	}
	writeJSON(w, http.StatusOK, child)
}

func (h *ChildrenHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid id"})
		return
	}
	child, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "child not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		return
	}
	writeJSON(w, http.StatusOK, child)
}
