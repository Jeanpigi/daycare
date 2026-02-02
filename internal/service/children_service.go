package service

import (
	"context"
	"strings"

	"daycare/internal/domain"
	mysqlrepo "daycare/internal/repository/mysql"
)

type ChildrenService struct{ repo *mysqlrepo.ChildrenRepo }

func NewChildrenService(r *mysqlrepo.ChildrenRepo) *ChildrenService { return &ChildrenService{repo: r} }

type CreateChildInput struct {
	DocumentNumber string
	FirstName      string
	LastName       string
	GuardianName   string
	GuardianPhone  string
}

func (s *ChildrenService) Create(ctx context.Context, in CreateChildInput) (domain.Child, error) {
	in.DocumentNumber = strings.TrimSpace(in.DocumentNumber)
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)
	in.GuardianName = strings.TrimSpace(in.GuardianName)
	in.GuardianPhone = strings.TrimSpace(in.GuardianPhone)

	if in.DocumentNumber == "" || in.FirstName == "" || in.LastName == "" || in.GuardianName == "" {
		return domain.Child{}, domain.ErrInvalid
	}

	var phone *string
	if in.GuardianPhone != "" {
		phone = &in.GuardianPhone
	}
	id, err := s.repo.Create(ctx, domain.Child{
		DocumentNumber: in.DocumentNumber,
		FirstName:      in.FirstName,
		LastName:       in.LastName,
		GuardianName:   in.GuardianName,
		GuardianPhone:  phone,
	})
	if err != nil {
		return domain.Child{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ChildrenService) GetByDocument(ctx context.Context, doc string) (domain.Child, error) {
	return s.repo.GetByDocument(ctx, doc)
}

func (s *ChildrenService) GetByID(ctx context.Context, id uint64) (domain.Child, error) {
	return s.repo.GetByID(ctx, id)
}
