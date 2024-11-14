package services

import (
	"context"
	"fmt"

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

func (s *businessService) ListBySearch(ctx context.Context, search string) ([]entity.Business, error) {
	business, err := s.repos.Business.ListBySearch(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("failed to list businesses: %w", err)
	}

	return business, nil
}

func (s *businessService) ListServicesBySearch(ctx context.Context, search string) ([]entity.BusinessService, error) {
	services, err := s.repos.Service.ListServicesBySearch(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	return services, nil
}

func (s *businessService) Create(ctx context.Context, business *entity.Business) error {
	// Validate business data
	if business.Name == "" {
		return fmt.Errorf("business name is required")
	}

	// Create business
	if err := s.repos.Business.Create(ctx, business); err != nil {
		return fmt.Errorf("failed to create business: %w", err)
	}

	return nil
}

func (s *businessService) Get(ctx context.Context, id int) (*entity.Business, error) {
	business, err := s.repos.Business.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get business: %w", err)
	}

	return business, nil
}

func (s *businessService) Update(ctx context.Context, business *entity.Business) error {
	// Validate business existence
	existing, err := s.repos.Business.Get(ctx, business.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing business: %w", err)
	}

	// Validate business data
	if business.Name == "" {
		return fmt.Errorf("business name is required")
	}

	// Update only allows changing the name
	existing.Name = business.Name

	if err := s.repos.Business.Update(ctx, existing); err != nil {
		return fmt.Errorf("failed to update business: %w", err)
	}

	return nil
}

func (s *businessService) UpdateAppearance(ctx context.Context, id int, logoURL string, colorScheme map[string]interface{}) error {
	// Validate business existence
	if _, err := s.repos.Business.Get(ctx, id); err != nil {
		return fmt.Errorf("failed to get existing business: %w", err)
	}

	// Validate color scheme
	if colorScheme != nil {
		if err := validateColorScheme(colorScheme); err != nil {
			return fmt.Errorf("invalid color scheme: %w", err)
		}
	}

	if err := s.repos.Business.UpdateAppearance(ctx, id, logoURL, colorScheme); err != nil {
		return fmt.Errorf("failed to update business appearance: %w", err)
	}

	return nil
}

// validateColorScheme ensures the color scheme contains valid values
func validateColorScheme(colorScheme map[string]interface{}) error {
	requiredColors := []string{"primary", "secondary", "background"}
	for _, color := range requiredColors {
		if _, ok := colorScheme[color]; !ok {
			return fmt.Errorf("missing required color: %s", color)
		}
	}
	return nil
}
