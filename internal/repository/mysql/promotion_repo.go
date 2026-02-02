package mysqlrepo

import (
	"context"
	"database/sql"
	"time"

	"daycare/internal/domain"
)

type PromotionRepo struct{ db *sql.DB }

func NewPromotionRepo(db *sql.DB) *PromotionRepo { return &PromotionRepo{db: db} }

func (r *PromotionRepo) ListActive(ctx context.Context, now time.Time) ([]domain.Promotion, error) {
	// vigencia: (starts_at is null or starts_at <= now) AND (ends_at is null or ends_at >= now)
	q := `
		SELECT id, name, rule_type, promo_price, min_days, min_minutes, starts_at, ends_at, priority, active, created_by, created_at
		FROM promotions
		WHERE active=1
		  AND (starts_at IS NULL OR starts_at <= ?)
		  AND (ends_at IS NULL OR ends_at >= ?)
		ORDER BY priority ASC, id ASC
	`
	rows, err := r.db.QueryContext(ctx, q, now, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Promotion
	for rows.Next() {
		var p domain.Promotion
		var minDays, minMins sql.NullInt32
		var sAt, eAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.Name, &p.RuleType, &p.PromoPrice, &minDays, &minMins, &sAt, &eAt, &p.Priority, &p.Active, &p.CreatedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		if minDays.Valid {
			v := int(minDays.Int32)
			p.MinDays = &v
		}
		if minMins.Valid {
			v := int(minMins.Int32)
			p.MinMinutes = &v
		}
		if sAt.Valid {
			t := sAt.Time
			p.StartsAt = &t
		}
		if eAt.Valid {
			t := eAt.Time
			p.EndsAt = &t
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PromotionRepo) Create(ctx context.Context, p domain.Promotion, actorID uint64) (uint64, error) {
	q := `
		INSERT INTO promotions (name, rule_type, promo_price, min_days, min_minutes, starts_at, ends_at, priority, active, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	var minDays any = nil
	var minMins any = nil
	var sAt any = nil
	var eAt any = nil
	if p.MinDays != nil {
		minDays = *p.MinDays
	}
	if p.MinMinutes != nil {
		minMins = *p.MinMinutes
	}
	if p.StartsAt != nil {
		sAt = *p.StartsAt
	}
	if p.EndsAt != nil {
		eAt = *p.EndsAt
	}
	active := 1
	if !p.Active {
		active = 0
	}
	res, err := r.db.ExecContext(ctx, q, p.Name, p.RuleType, p.PromoPrice, minDays, minMins, sAt, eAt, p.Priority, active, actorID)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *PromotionRepo) SetActive(ctx context.Context, promoID uint64, active bool) error {
	val := 0
	if active {
		val = 1
	}
	_, err := r.db.ExecContext(ctx, `UPDATE promotions SET active=? WHERE id=?`, val, promoID)
	return err
}
