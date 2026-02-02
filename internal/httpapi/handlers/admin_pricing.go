package handlers

import (
	"encoding/json"
	"net/http"

	"daycare/internal/domain"
	"daycare/internal/httpapi/middleware"
	mysqlrepo "daycare/internal/repository/mysql"
	"daycare/internal/service"
)

type AdminPricingHandler struct {
	admin  *service.AdminService
	prRepo *mysqlrepo.PricingRepo
}

func NewAdminPricingHandler(admin *service.AdminService, pr *mysqlrepo.PricingRepo) *AdminPricingHandler {
	return &AdminPricingHandler{admin: admin, prRepo: pr}
}

func (h *AdminPricingHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	p, err := h.prRepo.GetActive(r.Context())
	if err != nil {
		if err == domain.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "no active pricing"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

type createPricingReq struct {
	StandardPrice float64 `json:"standard_price"`
	Currency      string  `json:"currency"`
}

func (h *AdminPricingHandler) CreateAndActivate(w http.ResponseWriter, r *http.Request) {
	var req createPricingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	if req.StandardPrice <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "standard_price must be > 0"})
		return
	}
	if req.Currency == "" {
		req.Currency = "COP"
	}

	actorID := middleware.UserID(r.Context())
	id, err := h.admin.CreateAndActivatePricing(r.Context(), req.StandardPrice, req.Currency, actorID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"active_pricing_id": id})
}
