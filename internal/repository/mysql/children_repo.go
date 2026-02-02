package mysqlrepo

import (
	"context"
	"database/sql"
	"strings"

	"daycare/internal/domain"
)

type ChildrenRepo struct{ db *sql.DB }

func NewChildrenRepo(db *sql.DB) *ChildrenRepo { return &ChildrenRepo{db: db} }

func normalizeDocument(doc string) string {
	doc = strings.TrimSpace(doc)
	doc = strings.ReplaceAll(doc, ".", "")
	doc = strings.ReplaceAll(doc, " ", "")
	return doc
}

func (r *ChildrenRepo) Create(ctx context.Context, c domain.Child) (uint64, error) {
	c.DocumentNumber = normalizeDocument(c.DocumentNumber)

	q := `INSERT INTO children (document_number, first_name, last_name, guardian_name, guardian_phone)
	      VALUES (?, ?, ?, ?, ?)`
	var phone any = nil
	if c.GuardianPhone != nil && strings.TrimSpace(*c.GuardianPhone) != "" {
		phone = strings.TrimSpace(*c.GuardianPhone)
	}
	res, err := r.db.ExecContext(ctx, q, c.DocumentNumber, c.FirstName, c.LastName, c.GuardianName, phone)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return 0, domain.ErrConflict
		}
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *ChildrenRepo) GetByDocument(ctx context.Context, doc string) (domain.Child, error) {
	doc = normalizeDocument(doc)
	q := `SELECT id, document_number, first_name, last_name, guardian_name, guardian_phone, created_at
	      FROM children WHERE document_number=?`
	var c domain.Child
	var phone sql.NullString
	err := r.db.QueryRowContext(ctx, q, doc).Scan(&c.ID, &c.DocumentNumber, &c.FirstName, &c.LastName, &c.GuardianName, &phone, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.Child{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Child{}, err
	}
	if phone.Valid {
		c.GuardianPhone = &phone.String
	}
	return c, nil
}

func (r *ChildrenRepo) GetByID(ctx context.Context, id uint64) (domain.Child, error) {
	q := `SELECT id, document_number, first_name, last_name, guardian_name, guardian_phone, created_at
	      FROM children WHERE id=?`
	var c domain.Child
	var phone sql.NullString
	err := r.db.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.DocumentNumber, &c.FirstName, &c.LastName, &c.GuardianName, &phone, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.Child{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Child{}, err
	}
	if phone.Valid {
		c.GuardianPhone = &phone.String
	}
	return c, nil
}
