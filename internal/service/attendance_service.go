package service

import (
	"context"
	"time"

	"daycare/internal/domain"
	mysqlrepo "daycare/internal/repository/mysql"
)

type AttendanceService struct {
	children *mysqlrepo.ChildrenRepo
	att      *mysqlrepo.AttendanceRepo
	pricing  *mysqlrepo.PricingRepo
	promos   *mysqlrepo.PromotionRepo
	calc     *PricingService
}

func NewAttendanceService(
	children *mysqlrepo.ChildrenRepo,
	att *mysqlrepo.AttendanceRepo,
	pricing *mysqlrepo.PricingRepo,
	promos *mysqlrepo.PromotionRepo,
	calc *PricingService,
) *AttendanceService {
	return &AttendanceService{children: children, att: att, pricing: pricing, promos: promos, calc: calc}
}

type CheckInResult struct {
	AttendanceID uint64
	Child        domain.Child
	CheckedInAt  time.Time
}

func monthRange(t time.Time) (time.Time, time.Time) {
	loc := t.Location()
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)
	return start, end
}

func (s *AttendanceService) CheckInByDocument(ctx context.Context, doc string, actorID uint64) (CheckInResult, error) {
	child, err := s.children.GetByDocument(ctx, doc)
	if err != nil {
		return CheckInResult{}, err
	}

	now := time.Now()
	id, err := s.att.CreateCheckIn(ctx, child.ID, now, actorID)
	if err != nil {
		return CheckInResult{}, err
	}
	return CheckInResult{AttendanceID: id, Child: child, CheckedInAt: now}, nil
}

type CheckOutResult struct {
	AttendanceID uint64
	Child        domain.Child
	CheckedInAt  time.Time
	CheckedOutAt time.Time
	Minutes      int
	Pricing      PricingResult
	MonthlyDays  int
	MonthlyMins  int
}

func (s *AttendanceService) CheckOutByDocument(ctx context.Context, doc string, actorID uint64) (CheckOutResult, error) {
	child, err := s.children.GetByDocument(ctx, doc)
	if err != nil {
		return CheckOutResult{}, err
	}

	open, err := s.att.GetOpenByChild(ctx, child.ID)
	if err != nil {
		return CheckOutResult{}, err
	}

	now := time.Now()
	mins := int(now.Sub(open.CheckedInAt).Minutes())
	if mins < 0 {
		mins = 0
	}

	// pricing activo + promos activas
	prc, err := s.pricing.GetActive(ctx)
	if err != nil {
		return CheckOutResult{}, err
	}
	activePromos, err := s.promos.ListActive(ctx, now)
	if err != nil {
		return CheckOutResult{}, err
	}

	// stats mes (asistencias cerradas) + sumar esta jornada para que promo aplique inmediatamente
	ms, me := monthRange(now)
	stats, err := s.att.GetMonthlyStatsClosed(ctx, child.ID, ms, me)
	if err != nil {
		return CheckOutResult{}, err
	}
	stats.TotalMinutes += mins
	stats.DaysVisited += 1

	price := s.calc.Calculate(prc, activePromos, stats, now)

	var promoID *uint64 = nil
	if price.PromoID != nil {
		promoID = price.PromoID
	}

	prID := prc.ID
	if err := s.att.Close(ctx, open.ID, now, mins, &prID, promoID, price.Gross, price.Discount, price.Net, price.Currency, actorID); err != nil {
		return CheckOutResult{}, err
	}

	return CheckOutResult{
		AttendanceID: open.ID,
		Child:        child,
		CheckedInAt:  open.CheckedInAt,
		CheckedOutAt: now,
		Minutes:      mins,
		Pricing:      price,
		MonthlyDays:  stats.DaysVisited,
		MonthlyMins:  stats.TotalMinutes,
	}, nil
}
