package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"daycare/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

var ErrNotFound = errors.New("not found")

// Cuenta admins existentes
func (r *UserRepo) CountAdmins(ctx context.Context) (int, error) {
	const q = `SELECT COUNT(*) FROM users WHERE role='ADMIN'`
	var n int
	if err := r.db.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

// Crea usuario (ADMIN/STAFF). Email es UNIQUE en DB.
func (r *UserRepo) Create(ctx context.Context, u domain.User) (uint64, error) {
	const q = `INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, q, u.Name, u.Email, u.PasswordHash, u.Role)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	const q = `
SELECT id, name, email, password_hash, role, created_at
FROM users
WHERE email = ?
LIMIT 1`

	var u domain.User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}
