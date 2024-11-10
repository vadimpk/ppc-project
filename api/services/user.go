package services

import (
	"context"
	"errors"
	"fmt"

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

func (s *userService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	// Validate business existence if not admin
	if user.Role != entity.RoleAdmin {
		business, err := s.repos.Business.Get(ctx, user.BusinessID)
		if err != nil {
			return nil, fmt.Errorf("invalid business: %w", err)
		}
		user.BusinessID = business.ID
	}

	// Check unique constraints
	if user.Email != nil {
		_, err := s.repos.User.GetByEmail(ctx, user.BusinessID, *user.Email)
		if err == nil {
			return nil, fmt.Errorf("email already exists")
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
	}

	if user.Phone != nil {
		_, err := s.repos.User.GetByPhone(ctx, user.BusinessID, *user.Phone)
		if err == nil {
			return nil, fmt.Errorf("phone already exists")
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("failed to check phone uniqueness: %w", err)
		}
	}

	if err := s.repos.User.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) CreateBusinessAdmin(ctx context.Context, businessName string, user *entity.User) (*entity.User, error) {
	// Check unique constraints for email/phone
	if user.Email != nil {
		exists, err := s.repos.User.GetByEmail(ctx, user.BusinessID, *user.Email)
		if err == nil && exists != nil {
			return nil, fmt.Errorf("email already exists")
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
	}

	if user.Phone != nil {
		exists, err := s.repos.User.GetByPhone(ctx, user.BusinessID, *user.Phone)
		if err == nil && exists != nil {
			return nil, fmt.Errorf("phone already exists")
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("failed to check phone uniqueness: %w", err)
		}
	}

	if err := s.repos.User.CreateBusinessAdmin(ctx, businessName, user); err != nil {
		return nil, fmt.Errorf("failed to create business admin: %w", err)
	}

	return user, nil
}

func (s *userService) Get(ctx context.Context, id int) (*entity.User, error) {
	user, err := s.repos.User.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, businessID int, email string) (*entity.User, error) {
	user, err := s.repos.User.GetByEmail(ctx, businessID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *userService) GetByPhone(ctx context.Context, businessID int, phone string) (*entity.User, error) {
	user, err := s.repos.User.GetByPhone(ctx, businessID, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	// Check if user exists
	existing, err := s.repos.User.Get(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user: %w", err)
	}

	// Check unique constraints if email/phone is being updated
	if user.Email != nil && (existing.Email == nil || *existing.Email != *user.Email) {
		_, err := s.repos.User.GetByEmail(ctx, user.BusinessID, *user.Email)
		if err == nil {
			return nil, fmt.Errorf("email already exists")
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
	}

	if user.Phone != nil && (existing.Phone == nil || *existing.Phone != *user.Phone) {
		_, err := s.repos.User.GetByPhone(ctx, user.BusinessID, *user.Phone)
		if err == nil {
			return nil, fmt.Errorf("phone already exists")
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("failed to check phone uniqueness: %w", err)
		}
	}

	if err := s.repos.User.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) Authenticate(ctx context.Context, businessID int, email string, phone string, password string) (*entity.User, error) {
	var user *entity.User
	var err error

	if email != "" {
		user, err = s.repos.User.GetByEmail(ctx, businessID, email)
	} else if phone != "" {
		user, err = s.repos.User.GetByPhone(ctx, businessID, phone)
	} else {
		return nil, fmt.Errorf("either email or phone is required")
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Password comparison is handled at the handler level
	return user, nil
}
