package services

import (
	"context"
	"fmt"

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
	// Validate business existence
	if _, err := s.repos.Business.Get(ctx, service.BusinessID); err != nil {
		return fmt.Errorf("invalid business: %w", err)
	}

	// Validate service data
	if err := validateServiceData(service); err != nil {
		return err
	}

	// Create service
	if err := s.repos.Service.Create(ctx, service); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	return nil
}

func (s *businessServiceService) Get(ctx context.Context, id int) (*entity.BusinessService, error) {
	service, err := s.repos.Service.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return service, nil
}

func (s *businessServiceService) Update(ctx context.Context, service *entity.BusinessService) error {
	// Validate service existence
	existing, err := s.repos.Service.Get(ctx, service.ID)
	if err != nil {
		return fmt.Errorf("invalid service: %w", err)
	}

	// Ensure business ID matches
	if existing.BusinessID != service.BusinessID {
		return fmt.Errorf("service does not belong to the business")
	}

	// Validate service data
	if err := validateServiceData(service); err != nil {
		return err
	}

	// Update service
	if err := s.repos.Service.Update(ctx, service); err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	return nil
}

func (s *businessServiceService) Delete(ctx context.Context, id int) error {
	// Validate service existence
	if _, err := s.repos.Service.Get(ctx, id); err != nil {
		return fmt.Errorf("invalid service: %w", err)
	}

	// Soft delete the service
	if err := s.repos.Service.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

func (s *businessServiceService) List(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	// Validate business existence
	if _, err := s.repos.Business.Get(ctx, businessID); err != nil {
		return nil, fmt.Errorf("invalid business: %w", err)
	}

	services, err := s.repos.Service.List(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	return services, nil
}

func (s *businessServiceService) ListActive(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	// Validate business existence
	if _, err := s.repos.Business.Get(ctx, businessID); err != nil {
		return nil, fmt.Errorf("invalid business: %w", err)
	}

	services, err := s.repos.Service.ListActive(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to list active services: %w", err)
	}

	return services, nil
}

func validateServiceData(service *entity.BusinessService) error {
	if service.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if service.Duration <= 0 {
		return fmt.Errorf("service duration must be positive")
	}
	if service.Price < 0 {
		return fmt.Errorf("service price cannot be negative")
	}
	return nil
}
