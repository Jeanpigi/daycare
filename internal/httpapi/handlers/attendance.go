package handlers

import (
	"encoding/json"
	"net/http"

	"daycare/internal/domain"
	"daycare/internal/httpapi/middleware"
	"daycare/internal/service"
)

type AttendanceHandler struct {
	svc *service.AttendanceService
}

func NewAttendanceHandler(s *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: s}
}

type docReq struct {
	DocumentNumber string `json:"document_number"`
}

func (h *AttendanceHandler) CheckInByDocument(w http.ResponseWriter, r *http.Request) {
	var req docReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	if req.DocumentNumber == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "document_number required"})
		return
	}

	actorID := middleware.UserID(r.Context())
	res, err := h.svc.CheckInByDocument(r.Context(), req.DocumentNumber, actorID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "child not found"})
		case domain.ErrConflict:
			writeJSON(w, http.StatusConflict, map[string]any{"error": "child already checked-in"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"attendance_id":   res.AttendanceID,
		"child_id":        res.Child.ID,
		"document_number": res.Child.DocumentNumber,
		"checked_in_at":   res.CheckedInAt,
	})
}

func (h *AttendanceHandler) CheckOutByDocument(w http.ResponseWriter, r *http.Request) {
	var req docReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	if req.DocumentNumber == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "document_number required"})
		return
	}

	actorID := middleware.UserID(r.Context())
	res, err := h.svc.CheckOutByDocument(r.Context(), req.DocumentNumber, actorID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "child or open attendance not found"})
		case domain.ErrConflict:
			writeJSON(w, http.StatusConflict, map[string]any{"error": "no open attendance to close"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "server error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"attendance_id":   res.AttendanceID,
		"child_id":        res.Child.ID,
		"document_number": res.Child.DocumentNumber,
		"checked_in_at":   res.CheckedInAt,
		"checked_out_at":  res.CheckedOutAt,
		"minutes":         res.Minutes,

		"pricing": map[string]any{
			"pricing_id":      res.Pricing.PricingID,
			"promo_id":        res.Pricing.PromoID,
			"gross_amount":    res.Pricing.Gross,
			"discount_amount": res.Pricing.Discount,
			"net_amount":      res.Pricing.Net,
			"currency":        res.Pricing.Currency,
			"used_promo":      res.Pricing.UsedPromo,
			"promo_name":      res.Pricing.PromoName,
		},
		"monthly_days":    res.MonthlyDays,
		"monthly_minutes": res.MonthlyMins,
	})
}
