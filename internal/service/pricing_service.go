package service

import (
	"time"

	"daycare/internal/domain"
)

type PricingResult struct {
	PricingID uint64
	PromoID   *uint64
	Gross     float64
	Discount  float64
	Net       float64
	Currency  string
	UsedPromo bool
	PromoName *string
	PromoRule *string
}

type PricingService struct{}

func NewPricingService() *PricingService { return &PricingService{} }

// Evalúa promos activas. Por ahora implementamos regla:
// - LOYALTY_MONTH: aplica si days>=min_days OR minutes>=min_minutes
func (s *PricingService) Calculate(pricing domain.PricingConfig, promos []domain.Promotion, stats domain.MonthlyStats, now time.Time) PricingResult {
	res := PricingResult{
		PricingID: pricing.ID,
		Gross:     pricing.StandardPrice,
		Discount:  0,
		Net:       pricing.StandardPrice,
		Currency:  pricing.Currency,
		UsedPromo: false,
	}

	for _, p := range promos {
		if p.RuleType != "LOYALTY_MONTH" {
			continue
		}

		apply := false
		if p.MinDays != nil && stats.DaysVisited >= *p.MinDays {
			apply = true
		}
		if p.MinMinutes != nil && stats.TotalMinutes >= *p.MinMinutes {
			apply = true
		}

		if apply {
			res.UsedPromo = true
			res.Net = p.PromoPrice
			disc := pricing.StandardPrice - p.PromoPrice
			if disc < 0 {
				disc = 0
			}
			res.Discount = disc
			id := p.ID
			res.PromoID = &id
			name := p.Name
			rule := p.RuleType
			res.PromoName = &name
			res.PromoRule = &rule
			return res // primera por prioridad (promos vienen ordenadas)
		}
	}

	return res
}
