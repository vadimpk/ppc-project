package repository

import (
	"context"
	"database/sql"

	"github.com/vadimpk/ppc-project/entity"
)

type BusinessServiceRepository interface {
	Create(ctx context.Context, service *entity.BusinessService) error
	Get(ctx context.Context, id int) (*entity.BusinessService, error)
	Update(ctx context.Context, service *entity.BusinessService) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, businessID int) ([]entity.BusinessService, error)
	ListActive(ctx context.Context, businessID int) ([]entity.BusinessService, error)
}

type businessServiceRepository struct {
	db *sql.DB
}

func NewBusinessServiceRepository(db *sql.DB) BusinessServiceRepository {
	return &businessServiceRepository{
		db: db,
	}
}

func (r *businessServiceRepository) Create(ctx context.Context, service *entity.BusinessService) error {
	panic("not implemented")
}

func (r *businessServiceRepository) Get(ctx context.Context, id int) (*entity.BusinessService, error) {
	panic("not implemented")
}

func (r *businessServiceRepository) Update(ctx context.Context, service *entity.BusinessService) error {
	panic("not implemented")
}

func (r *businessServiceRepository) Delete(ctx context.Context, id int) error {
	panic("not implemented")
}

func (r *businessServiceRepository) List(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	panic("not implemented")
}

func (r *businessServiceRepository) ListActive(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	panic("not implemented")
}
