package service

import (
	"context"
	"errors"
	"strings"

	"daycare/internal/domain"
	mysqlrepo "daycare/internal/repository/mysql"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
)

type AdminUsersService struct {
	users *mysqlrepo.UserRepo
}

func NewAdminUsersService(users *mysqlrepo.UserRepo) *AdminUsersService {
	return &AdminUsersService{users: users}
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	Role     string // "ADMIN" o "STAFF"
}

func (s *AdminUsersService) CreateUser(ctx context.Context, in CreateUserInput) (uint64, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	pass := strings.TrimSpace(in.Password)
	role := strings.ToUpper(strings.TrimSpace(in.Role))
	if role == "" {
		role = "STAFF"
	}

	if name == "" || email == "" || len(pass) < 6 {
		return 0, ErrInvalidInput
	}
	if role != "ADMIN" && role != "STAFF" {
		return 0, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	id, err := s.users.Create(ctx, domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	})
	if err != nil {
		// MySQL duplicate key
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return 0, ErrConflict
		}
		return 0, err
	}

	return id, nil
}
