package repository

import (
	"context"
	"database/sql"

	"github.com/vadimpk/ppc-project/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Get(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, businessID int, email string) (*entity.User, error)
	GetByPhone(ctx context.Context, businessID int, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	panic("not implemented")
}

func (r *userRepository) Get(ctx context.Context, id int) (*entity.User, error) {
	panic("not implemented")
}

func (r *userRepository) GetByEmail(ctx context.Context, businessID int, email string) (*entity.User, error) {
	panic("not implemented")
}

func (r *userRepository) GetByPhone(ctx context.Context, businessID int, phone string) (*entity.User, error) {
	panic("not implemented")
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	panic("not implemented")
}
