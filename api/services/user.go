package services

import (
	"context"

	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

type userService struct {
	repos *repository.Repositories
}

func NewUserService(repos *repository.Repositories) UserService {
	return &userService{
		repos: repos,
	}
}

func (s *userService) Create(ctx context.Context, user *entity.User) error {
	panic("not implemented")
}

func (s *userService) Get(ctx context.Context, id int) (*entity.User, error) {
	panic("not implemented")
}

func (s *userService) GetByEmail(ctx context.Context, businessID int, email string) (*entity.User, error) {
	panic("not implemented")
}

func (s *userService) GetByPhone(ctx context.Context, businessID int, phone string) (*entity.User, error) {
	panic("not implemented")
}

func (s *userService) Update(ctx context.Context, user *entity.User) error {
	panic("not implemented")
}

func (s *userService) Authenticate(ctx context.Context, businessID int, email string, phone string, password string) (*entity.User, error) {
	panic("not implemented")
}
