package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository/db/sqlc"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --dir . --name BusinessServiceRepository --output ./mocks
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
		return r.db.HandleBasicErrors(err)
	}

	service.ID = int(dbService.ID)
	service.CreatedAt = dbService.CreatedAt.Time
	return nil
}

func (r *businessServiceRepository) Get(ctx context.Context, id int) (*entity.BusinessService, error) {
	dbService, err := r.db.SQLC.GetService(ctx, int32(id))
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
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
		return r.db.HandleBasicErrors(err)
	}

	service.CreatedAt = dbService.CreatedAt.Time
	return nil
}

func (r *businessServiceRepository) Delete(ctx context.Context, id int) error {
	err := r.db.SQLC.DeleteService(ctx, int32(id))
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	return nil
}

func (r *businessServiceRepository) List(ctx context.Context, businessID int) ([]entity.BusinessService, error) {
	dbServices, err := r.db.SQLC.ListServices(ctx, pgtype.Int4{Int32: int32(businessID), Valid: true})
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
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
		return nil, r.db.HandleBasicErrors(err)
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
