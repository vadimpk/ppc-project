package services

import (
	"context"

	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

type employeeService struct {
	repos *repository.Repositories
}

func NewEmployeeService(repos *repository.Repositories) EmployeeService {
	return &employeeService{
		repos: repos,
	}
}

func (s *employeeService) Create(ctx context.Context, employee *entity.Employee) error {
	panic("not implemented")
}

func (s *employeeService) Get(ctx context.Context, id int) (*entity.Employee, error) {
	panic("not implemented")
}

func (s *employeeService) Update(ctx context.Context, employee *entity.Employee) error {
	panic("not implemented")
}

func (s *employeeService) List(ctx context.Context, businessID int) ([]entity.Employee, error) {
	panic("not implemented")
}

func (s *employeeService) AssignServices(ctx context.Context, employeeID int, serviceIDs []int) error {
	panic("not implemented")
}

func (s *employeeService) RemoveServices(ctx context.Context, employeeID int, serviceIDs []int) error {
	panic("not implemented")
}

func (s *employeeService) GetServices(ctx context.Context, employeeID int) ([]entity.BusinessService, error) {
	panic("not implemented")
}
