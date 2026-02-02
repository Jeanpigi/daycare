package domain

import "time"

type User struct {
	ID           uint64
	Name         string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

type Child struct {
	ID             uint64
	DocumentNumber string
	FirstName      string
	LastName       string
	GuardianName   string
	GuardianPhone  *string
	CreatedAt      time.Time
}

type PricingConfig struct {
	ID            uint64
	StandardPrice float64
	Currency      string
	Active        bool
	CreatedBy     uint64
	CreatedAt     time.Time
}

type Promotion struct {
	ID         uint64
	Name       string
	RuleType   string // LOYALTY_MONTH
	PromoPrice float64
	MinDays    *int
	MinMinutes *int
	StartsAt   *time.Time
	EndsAt     *time.Time
	Priority   int
	Active     bool
	CreatedBy  uint64
	CreatedAt  time.Time
}

type Attendance struct {
	ID           uint64
	ChildID      uint64
	CheckedInAt  time.Time
	CheckedOutAt *time.Time
	Minutes      *int

	PricingConfigID *uint64
	PromoID         *uint64

	GrossAmount    *float64
	DiscountAmount *float64
	NetAmount      *float64
	Currency       string

	CreatedAt time.Time
}

type MonthlyStats struct {
	DaysVisited  int
	TotalMinutes int
}
