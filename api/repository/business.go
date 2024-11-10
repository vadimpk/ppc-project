package repository

import (
	"context"

	"github.com/vadimpk/ppc-project/entity"
)

type BusinessRepository interface {
	Create(ctx context.Context, business *entity.Business) error
	Get(ctx context.Context, id int) (*entity.Business, error)
	Update(ctx context.Context, business *entity.Business) error
	UpdateAppearance(ctx context.Context, id int, logoURL string, colorScheme map[string]interface{}) error
}

type businessRepository struct {
	db *DB
}

func NewBusinessRepository(db *DB) BusinessRepository {
	return &businessRepository{
		db: db,
	}
}

func (r *businessRepository) Create(ctx context.Context, business *entity.Business) error {
	panic("not implemented")
}

func (r *businessRepository) Get(ctx context.Context, id int) (*entity.Business, error) {
	panic("not implemented")
}

func (r *businessRepository) Update(ctx context.Context, business *entity.Business) error {
	panic("not implemented")
}

func (r *businessRepository) UpdateAppearance(ctx context.Context, id int, logoURL string, colorScheme map[string]interface{}) error {
	panic("not implemented")
}
