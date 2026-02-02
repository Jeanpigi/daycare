package mysqlrepo

import (
	"context"
	"database/sql"
	"time"

	"daycare/internal/domain"
)

type AttendanceRepo struct{ db *sql.DB }

func NewAttendanceRepo(db *sql.DB) *AttendanceRepo { return &AttendanceRepo{db: db} }

func (r *AttendanceRepo) GetOpenByChild(ctx context.Context, childID uint64) (domain.Attendance, error) {
	q := `
		SELECT id, child_id, checked_in_at, checked_out_at, minutes, pricing_config_id, promo_id,
		       gross_amount, discount_amount, net_amount, currency, created_at
		FROM attendances
		WHERE child_id=? AND checked_out_at IS NULL
		ORDER BY checked_in_at DESC
		LIMIT 1
	`
	var a domain.Attendance
	var outAt sql.NullTime
	var mins sql.NullInt32
	var pricingID sql.NullInt64
	var promoID sql.NullInt64
	var gross, disc, net sql.NullFloat64
	err := r.db.QueryRowContext(ctx, q, childID).Scan(
		&a.ID, &a.ChildID, &a.CheckedInAt, &outAt, &mins, &pricingID, &promoID,
		&gross, &disc, &net, &a.Currency, &a.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return domain.Attendance{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Attendance{}, err
	}
	if outAt.Valid {
		t := outAt.Time
		a.CheckedOutAt = &t
	}
	if mins.Valid {
		v := int(mins.Int32)
		a.Minutes = &v
	}
	if pricingID.Valid {
		v := uint64(pricingID.Int64)
		a.PricingConfigID = &v
	}
	if promoID.Valid {
		v := uint64(promoID.Int64)
		a.PromoID = &v
	}
	if gross.Valid {
		v := gross.Float64
		a.GrossAmount = &v
	}
	if disc.Valid {
		v := disc.Float64
		a.DiscountAmount = &v
	}
	if net.Valid {
		v := net.Float64
		a.NetAmount = &v
	}
	return a, nil
}

func (r *AttendanceRepo) CreateCheckIn(ctx context.Context, childID uint64, checkedInAt time.Time, createdBy uint64) (uint64, error) {
	// evitar doble open (regla en app)
	if _, err := r.GetOpenByChild(ctx, childID); err == nil {
		return 0, domain.ErrConflict
	} else if err != nil && err != domain.ErrNotFound {
		return 0, err
	}

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO attendances (child_id, checked_in_at, created_by) VALUES (?, ?, ?)`,
		childID, checkedInAt, createdBy,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *AttendanceRepo) Close(ctx context.Context, attendanceID uint64, checkedOutAt time.Time, minutes int,
	pricingID *uint64, promoID *uint64, gross, discount, net float64, currency string, closedBy uint64,
) error {
	q := `
		UPDATE attendances
		SET checked_out_at=?, minutes=?,
		    pricing_config_id=?, promo_id=?,
		    gross_amount=?, discount_amount=?, net_amount=?, currency=?,
		    closed_by=?
		WHERE id=? AND checked_out_at IS NULL
	`
	var pID any = nil
	var prID any = nil
	if pricingID != nil {
		prID = *pricingID
	}
	if promoID != nil {
		pID = *promoID
	}
	res, err := r.db.ExecContext(ctx, q, checkedOutAt, minutes, prID, pID, gross, discount, net, currency, closedBy, attendanceID)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return domain.ErrConflict
	}
	return nil
}

func (r *AttendanceRepo) GetMonthlyStatsClosed(ctx context.Context, childID uint64, from, to time.Time) (domain.MonthlyStats, error) {
	// días distintos de check_in en asistencias cerradas
	qDays := `
		SELECT COUNT(DISTINCT DATE(checked_in_at))
		FROM attendances
		WHERE child_id=? AND checked_out_at IS NOT NULL
		  AND checked_in_at >= ? AND checked_in_at < ?
	`
	var days int
	if err := r.db.QueryRowContext(ctx, qDays, childID, from, to).Scan(&days); err != nil {
		return domain.MonthlyStats{}, err
	}
	qMins := `
		SELECT COALESCE(SUM(minutes), 0)
		FROM attendances
		WHERE child_id=? AND checked_out_at IS NOT NULL
		  AND checked_in_at >= ? AND checked_in_at < ?
	`
	var mins int
	if err := r.db.QueryRowContext(ctx, qMins, childID, from, to).Scan(&mins); err != nil {
		return domain.MonthlyStats{}, err
	}
	return domain.MonthlyStats{DaysVisited: days, TotalMinutes: mins}, nil
}
