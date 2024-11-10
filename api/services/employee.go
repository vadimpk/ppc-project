package services

import (
	"context"
	"fmt"

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
	// Validate business existence
	if _, err := s.repos.Business.Get(ctx, employee.BusinessID); err != nil {
		return fmt.Errorf("invalid business: %w", err)
	}

	// Create user account for employee if not exists
	user, err := s.repos.User.Get(ctx, employee.UserID)
	if err != nil {
		return fmt.Errorf("invalid user: %w", err)
	}

	// Verify user belongs to the same business
	if user.BusinessID != employee.BusinessID {
		return fmt.Errorf("user does not belong to the business")
	}

	// Update user role to employee if it's not already
	if user.Role != entity.RoleEmployee {
		user.Role = entity.RoleEmployee
		if err := s.repos.User.Update(ctx, user); err != nil {
			return fmt.Errorf("failed to update user role: %w", err)
		}
	}

	// Create employee record
	if err := s.repos.Employee.Create(ctx, employee); err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}

	return nil
}

func (s *employeeService) Get(ctx context.Context, id int) (*entity.Employee, error) {
	employee, err := s.repos.Employee.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return employee, nil
}

func (s *employeeService) Update(ctx context.Context, employee *entity.Employee) error {
	// Validate employee existence and get current data
	existing, err := s.repos.Employee.Get(ctx, employee.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing employee: %w", err)
	}

	// Update only allowed fields
	existing.Specialization = employee.Specialization
	existing.IsActive = employee.IsActive

	if err := s.repos.Employee.Update(ctx, existing); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (s *employeeService) List(ctx context.Context, businessID int) ([]entity.Employee, error) {
	// Validate business existence
	if _, err := s.repos.Business.Get(ctx, businessID); err != nil {
		return nil, fmt.Errorf("invalid business: %w", err)
	}

	employees, err := s.repos.Employee.List(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to list employees: %w", err)
	}

	return employees, nil
}

func (s *employeeService) AssignServices(ctx context.Context, employeeID int, serviceIDs []int) error {
	// Validate employee existence and get current data
	employee, err := s.repos.Employee.Get(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("invalid employee: %w", err)
	}

	// Validate all services exist and belong to the same business
	for _, serviceID := range serviceIDs {
		service, err := s.repos.Service.Get(ctx, serviceID)
		if err != nil {
			return fmt.Errorf("invalid service ID %d: %w", serviceID, err)
		}
		if service.BusinessID != employee.BusinessID {
			return fmt.Errorf("service %d does not belong to the employee's business", serviceID)
		}
	}

	if err := s.repos.Employee.AssignServices(ctx, employeeID, serviceIDs); err != nil {
		return fmt.Errorf("failed to assign services: %w", err)
	}

	return nil
}

func (s *employeeService) RemoveServices(ctx context.Context, employeeID int, serviceIDs []int) error {
	// Validate employee existence
	if _, err := s.repos.Employee.Get(ctx, employeeID); err != nil {
		return fmt.Errorf("invalid employee: %w", err)
	}

	if err := s.repos.Employee.RemoveServices(ctx, employeeID, serviceIDs); err != nil {
		return fmt.Errorf("failed to remove services: %w", err)
	}

	return nil
}

func (s *employeeService) GetServices(ctx context.Context, employeeID int) ([]entity.BusinessService, error) {
	// Validate employee existence
	if _, err := s.repos.Employee.Get(ctx, employeeID); err != nil {
		return nil, fmt.Errorf("invalid employee: %w", err)
	}

	services, err := s.repos.Employee.GetServices(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee services: %w", err)
	}

	return services, nil
}
