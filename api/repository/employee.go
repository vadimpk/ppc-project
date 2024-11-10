package repository

import (
	"context"

	"github.com/vadimpk/ppc-project/entity"
)

type EmployeeRepository interface {
	Create(ctx context.Context, employee *entity.Employee) error
	Get(ctx context.Context, id int) (*entity.Employee, error)
	Update(ctx context.Context, employee *entity.Employee) error
	List(ctx context.Context, businessID int) ([]entity.Employee, error)
	AssignServices(ctx context.Context, employeeID int, serviceIDs []int) error
	RemoveServices(ctx context.Context, employeeID int, serviceIDs []int) error
	GetServices(ctx context.Context, employeeID int) ([]entity.BusinessService, error)
}

type employeeRepository struct {
	db *DB
}

func NewEmployeeRepository(db *DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (r *employeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
	panic("not implemented")
}

func (r *employeeRepository) Get(ctx context.Context, id int) (*entity.Employee, error) {
	panic("not implemented")
}

func (r *employeeRepository) Update(ctx context.Context, employee *entity.Employee) error {
	panic("not implemented")
}

func (r *employeeRepository) List(ctx context.Context, businessID int) ([]entity.Employee, error) {
	panic("not implemented")
}

func (r *employeeRepository) AssignServices(ctx context.Context, employeeID int, serviceIDs []int) error {
	panic("not implemented")
}

func (r *employeeRepository) RemoveServices(ctx context.Context, employeeID int, serviceIDs []int) error {
	panic("not implemented")
}

func (r *employeeRepository) GetServices(ctx context.Context, employeeID int) ([]entity.BusinessService, error) {
	panic("not implemented")
}
