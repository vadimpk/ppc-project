package services

import (
	"context"

	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

type businessServiceService struct {
	repos *repository.Repositories
}

func NewBusinessServiceService(repos *repository.Repositories) BusinessServiceService {
	return &businessServiceService{
		repos: repos,
	}
}

func (s *businessServiceService) Create(ctx context.Context, service *entity.BusinessService) error {
	panic("not implemented")
}

func (s *businessServiceService) Get(ctx context.Context, id int) (*entity.BusinessService, error) {
	panic("not implemented")
}

func (s *businessServiceService) Update(ctx context.Context, service *entity.BusinessService) error {
	panic("not implemented")
}

func (s *businessServiceService) Delete(ctx context.Context, id int) error {
	panic("not implemented")
}

func (s *businessServiceService) List(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	panic("not implemented")
}

func (s *businessServiceService) ListActive(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	panic("not implemented")
}
