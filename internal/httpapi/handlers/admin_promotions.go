package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"daycare/internal/domain"
	"daycare/internal/httpapi/middleware"
	"daycare/internal/service"

	"github.com/go-chi/chi/v5"
)

type AdminPromotionsHandler struct {
	admin *service.AdminService
}

func NewAdminPromotionsHandler(admin *service.AdminService) *AdminPromotionsHandler {
	return &AdminPromotionsHandler{admin: admin}
}

type createPromoReq struct {
	Name       string  `json:"name"`
	RuleType   string  `json:"rule_type"`   // LOYALTY_MONTH
	PromoPrice float64 `json:"promo_price"` // 15000

	MinDays    *int    `json:"min_days"`
	MinMinutes *int    `json:"min_minutes"`
	StartsAt   *string `json:"starts_at"` // RFC3339 opcional
	EndsAt     *string `json:"ends_at"`   // RFC3339 opcional

	Priority int  `json:"priority"`
	Active   bool `json:"active"`
}

func (h *AdminPromotionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPromoReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	if req.Name == "" || req.RuleType == "" || req.PromoPrice <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "name, rule_type, promo_price required"})
		return
	}
	if req.Priority == 0 {
		req.Priority = 100
	}

	var sAt *time.Time
	var eAt *time.Time
	if req.StartsAt != nil && *req.StartsAt != "" {
		t, err := time.Parse(time.RFC3339, *req.StartsAt)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "starts_at must be RFC3339"})
			return
		}
		sAt = &t
	}
	if req.EndsAt != nil && *req.EndsAt != "" {
		t, err := time.Parse(time.RFC3339, *req.EndsAt)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "ends_at must be RFC3339"})
			return
		}
		eAt = &t
	}

	actorID := middleware.UserID(r.Context())
	id, err := h.admin.CreatePromotion(r.Context(), domain.Promotion{
		Name:       req.Name,
		RuleType:   req.RuleType,
		PromoPrice: req.PromoPrice,
		MinDays:    req.MinDays,
		MinMinutes: req.MinMinutes,
		StartsAt:   sAt,
		EndsAt:     eAt,
		Priority:   req.Priority,
		Active:     req.Active,
	}, actorID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"promo_id": id})
}

func (h *AdminPromotionsHandler) Activate(w http.ResponseWriter, r *http.Request) {
	h.toggle(w, r, true)
}
func (h *AdminPromotionsHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	h.toggle(w, r, false)
}

func (h *AdminPromotionsHandler) toggle(w http.ResponseWriter, r *http.Request, active bool) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid id"})
		return
	}
	actorID := middleware.UserID(r.Context())
	if err := h.admin.SetPromotionActive(r.Context(), id, active, actorID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"promo_id": id, "active": active})
}
