package service

import (
	"context"
	"encoding/json"
	"time"

	"daycare/internal/domain"
	mysqlrepo "daycare/internal/repository/mysql"
)

type AdminService struct {
	pricing *mysqlrepo.PricingRepo
	promos  *mysqlrepo.PromotionRepo
	audit   *mysqlrepo.AuditRepo
}

func NewAdminService(pricing *mysqlrepo.PricingRepo, promos *mysqlrepo.PromotionRepo, audit *mysqlrepo.AuditRepo) *AdminService {
	return &AdminService{pricing: pricing, promos: promos, audit: audit}
}

func (s *AdminService) CreateAndActivatePricing(ctx context.Context, standard float64, currency string, actorID uint64) (uint64, error) {
	before := []byte{}
	if cur, err := s.pricing.GetActive(ctx); err == nil {
		b, _ := json.Marshal(cur)
		before = b
	}

	id, err := s.pricing.CreateAndActivate(ctx, standard, currency, actorID)
	if err != nil {
		return 0, err
	}

	afterObj, _ := s.pricing.GetActive(ctx)
	after, _ := json.Marshal(afterObj)
	_ = s.audit.Insert(ctx, actorID, "PRICING_UPDATE", "settings_pricing", &id, before, after)

	return id, nil
}

func (s *AdminService) CreatePromotion(ctx context.Context, p domain.Promotion, actorID uint64) (uint64, error) {
	id, err := s.promos.Create(ctx, p, actorID)
	if err != nil {
		return 0, err
	}
	after, _ := json.Marshal(map[string]any{
		"id": id, "name": p.Name, "rule_type": p.RuleType, "promo_price": p.PromoPrice,
	})
	_ = s.audit.Insert(ctx, actorID, "PROMO_CREATE", "promotions", &id, nil, after)
	return id, nil
}

func (s *AdminService) SetPromotionActive(ctx context.Context, promoID uint64, active bool, actorID uint64) error {
	before, _ := json.Marshal(map[string]any{"promo_id": promoID, "active": !active})
	if err := s.promos.SetActive(ctx, promoID, active); err != nil {
		return err
	}
	after, _ := json.Marshal(map[string]any{"promo_id": promoID, "active": active, "at": time.Now()})
	_ = s.audit.Insert(ctx, actorID, "PROMO_TOGGLE", "promotions", &promoID, before, after)
	return nil
}
