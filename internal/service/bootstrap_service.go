package service

import (
	"context"
	"strings"

	"daycare/internal/domain"
	mysqlrepo "daycare/internal/repository/mysql"

	"golang.org/x/crypto/bcrypt"
)

type BootstrapService struct {
	users   *mysqlrepo.UserRepo
	pricing *mysqlrepo.PricingRepo
	promos  *mysqlrepo.PromotionRepo
}

func NewBootstrapService(users *mysqlrepo.UserRepo, pricing *mysqlrepo.PricingRepo, promos *mysqlrepo.PromotionRepo) *BootstrapService {
	return &BootstrapService{users: users, pricing: pricing, promos: promos}
}

type BootstrapAdminInput struct {
	Name     string
	Email    string
	Password string
}

func (s *BootstrapService) CreateFirstAdmin(ctx context.Context, in BootstrapAdminInput) (uint64, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	pass := strings.TrimSpace(in.Password)

	if name == "" || email == "" || len(pass) < 6 {
		return 0, ErrInvalidInput
	}

	n, err := s.users.CountAdmins(ctx)
	if err != nil {
		return 0, err
	}
	if n > 0 {
		return 0, ErrConflict // ya existe admin
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	adminID, err := s.users.Create(ctx, domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "ADMIN",
	})
	if err != nil {
		return 0, err
	}

	// Defaults opcionales (para que el sistema tenga algo desde el día 1):
	_, _ = s.pricing.CreateAndActivate(ctx, 20000, "COP", adminID)

	minDays := 2
	minMins := 300
	_, _ = s.promos.Create(ctx, domain.Promotion{
		Name:       "Promo fidelidad",
		RuleType:   "LOYALTY_MONTH",
		PromoPrice: 15000,
		MinDays:    &minDays,
		MinMinutes: &minMins,
		Priority:   10,
		Active:     true,
	}, adminID)

	return adminID, nil
}
