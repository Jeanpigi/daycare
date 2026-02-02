package mysqlrepo

import (
	"context"
	"database/sql"

	"daycare/internal/domain"
)

type PricingRepo struct{ db *sql.DB }

func NewPricingRepo(db *sql.DB) *PricingRepo { return &PricingRepo{db: db} }

func (r *PricingRepo) GetActive(ctx context.Context) (domain.PricingConfig, error) {
	q := `SELECT id, standard_price, currency, active, created_by, created_at
	      FROM settings_pricing WHERE active=1 ORDER BY id DESC LIMIT 1`
	var p domain.PricingConfig
	err := r.db.QueryRowContext(ctx, q).Scan(&p.ID, &p.StandardPrice, &p.Currency, &p.Active, &p.CreatedBy, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.PricingConfig{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.PricingConfig{}, err
	}
	return p, nil
}

func (r *PricingRepo) CreateAndActivate(ctx context.Context, standard float64, currency string, actorID uint64) (uint64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	// desactivar actuales
	if _, err := tx.ExecContext(ctx, `UPDATE settings_pricing SET active=0 WHERE active=1`); err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx,
		`INSERT INTO settings_pricing (standard_price, currency, active, created_by) VALUES (?, ?, 1, ?)`,
		standard, currency, actorID,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return uint64(id), nil
}
