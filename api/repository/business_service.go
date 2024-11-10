package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository/db/sqlc"
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
	db *DB
}

func NewBusinessServiceRepository(db *DB) BusinessServiceRepository {
	return &businessServiceRepository{
		db: db,
	}
}

func (r *businessServiceRepository) Create(ctx context.Context, service *entity.BusinessService) error {
	var description pgtype.Text
	if service.Description != nil {
		description = r.db.ValidText(*service.Description)
	}

	dbService, err := r.db.SQLC.CreateService(ctx, sqlc.CreateServiceParams{
		BusinessID:  pgtype.Int4{Int32: int32(service.BusinessID), Valid: true},
		Name:        service.Name,
		Description: description,
		Duration:    int32(service.Duration),
		Price:       int32(service.Price),
		IsActive:    pgtype.Bool{Bool: service.IsActive, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	service.ID = int(dbService.ID)
	service.CreatedAt = dbService.CreatedAt.Time
	return nil
}

func (r *businessServiceRepository) Get(ctx context.Context, id int) (*entity.BusinessService, error) {
	dbService, err := r.db.SQLC.GetService(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return convertDBServiceToEntity(dbService), nil
}

func (r *businessServiceRepository) Update(ctx context.Context, service *entity.BusinessService) error {
	var description pgtype.Text
	if service.Description != nil {
		description = r.db.ValidText(*service.Description)
	}

	dbService, err := r.db.SQLC.UpdateService(ctx, sqlc.UpdateServiceParams{
		ID:          int32(service.ID),
		Name:        service.Name,
		Description: description,
		Duration:    int32(service.Duration),
		Price:       int32(service.Price),
		IsActive:    pgtype.Bool{Bool: service.IsActive, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	service.CreatedAt = dbService.CreatedAt.Time
	return nil
}

func (r *businessServiceRepository) Delete(ctx context.Context, id int) error {
	err := r.db.SQLC.DeleteService(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

func (r *businessServiceRepository) List(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	dbServices, err := r.db.SQLC.ListServices(ctx, pgtype.Int4{Int32: int32(businessID), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	services := make([]entity.BusinessService, len(dbServices))
	for i, dbService := range dbServices {
		services[i] = *convertDBServiceToEntity(dbService)
	}

	return services, nil
}

func (r *businessServiceRepository) ListActive(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	dbServices, err := r.db.SQLC.ListActiveServices(ctx, pgtype.Int4{Int32: int32(businessID), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list active services: %w", err)
	}

	services := make([]entity.BusinessService, len(dbServices))
	for i, dbService := range dbServices {
		services[i] = *convertDBServiceToEntity(dbService)
	}

	return services, nil
}

func convertDBServiceToEntity(s sqlc.Service) *entity.BusinessService {
	service := &entity.BusinessService{
		ID:         int(s.ID),
		BusinessID: int(s.BusinessID.Int32),
		Name:       s.Name,
		Duration:   int(s.Duration),
		Price:      int(s.Price),
		IsActive:   s.IsActive.Bool,
		CreatedAt:  s.CreatedAt.Time,
	}

	if s.Description.Valid {
		desc := s.Description.String
		service.Description = &desc
	}

	return service
}
