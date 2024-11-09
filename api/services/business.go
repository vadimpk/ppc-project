package services

import (
	"context"

	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

type businessService struct {
	repos *repository.Repositories
}

func NewBusinessService(repos *repository.Repositories) BusinessService {
	return &businessService{
		repos: repos,
	}
}

func (s *businessService) Create(ctx context.Context, business *entity.Business) error {
	panic("not implemented")
}

func (s *businessService) Get(ctx context.Context, id int) (*entity.Business, error) {
	panic("not implemented")
}

func (s *businessService) Update(ctx context.Context, business *entity.Business) error {
	panic("not implemented")
}

func (s *businessService) UpdateAppearance(ctx context.Context, id int, logoURL string, colorScheme map[string]interface{}) error {
	panic("not implemented")
}
