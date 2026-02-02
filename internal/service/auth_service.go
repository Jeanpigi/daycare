package service

import (
	"context"
	"time"

	"daycare/internal/domain"
	mysqlrepo "daycare/internal/repository/mysql"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users   *mysqlrepo.UserRepo
	secret  []byte
	ttlMins int
}

func NewAuthService(users *mysqlrepo.UserRepo, jwtSecret string, ttlMins int) *AuthService {
	return &AuthService{users: users, secret: []byte(jwtSecret), ttlMins: ttlMins}
}

type LoginResult struct {
	Token string
	User  domain.User
}

func (s *AuthService) Login(ctx context.Context, email, password string) (LoginResult, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return LoginResult{}, domain.ErrUnauthorized
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return LoginResult{}, domain.ErrUnauthorized
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"role": u.Role,
		"exp":  now.Add(time.Duration(s.ttlMins) * time.Minute).Unix(),
		"iat":  now.Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.secret)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: signed, User: u}, nil
}
