package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository/db/sqlc"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --dir . --name BusinessRepository --output ./mocks
type BusinessRepository interface {
	Create(ctx context.Context, business *entity.Business) error
	Get(ctx context.Context, id int) (*entity.Business, error)
	Update(ctx context.Context, business *entity.Business) error
	UpdateAppearance(ctx context.Context, id int, logoURL string, colorScheme map[string]interface{}) error
	ListBySearch(ctx context.Context, search string) ([]entity.Business, error)
}

type businessRepository struct {
	db *DB
}

func NewBusinessRepository(db *DB) BusinessRepository {
	return &businessRepository{
		db: db,
	}
}

func (r *businessRepository) ListBySearch(ctx context.Context, search string) ([]entity.Business, error) {
	dbBusinesses, err := r.db.SQLC.ListBySearch(ctx, "%"+search+"%")
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	businesses := make([]entity.Business, len(dbBusinesses))
	for i, dbBusiness := range dbBusinesses {
		businesses[i] = *convertDBBusinessToEntity(dbBusiness)
	}

	return businesses, nil
}

func (r *businessRepository) Create(ctx context.Context, business *entity.Business) error {
	var logoURL pgtype.Text
	if business.LogoURL != nil {
		logoURL = pgtype.Text{String: *business.LogoURL, Valid: true}
	}

	colorSchemeJSON, err := json.Marshal(business.ColorScheme)
	if err != nil {
		return fmt.Errorf("failed to marshal color scheme: %w", err)
	}

	dbBusiness, err := r.db.SQLC.CreateBusiness(ctx, sqlc.CreateBusinessParams{
		Name:        business.Name,
		LogoUrl:     logoURL,
		ColorScheme: colorSchemeJSON,
	})
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	business.ID = int(dbBusiness.ID)
	business.CreatedAt = dbBusiness.CreatedAt.Time
	return nil
}

func (r *businessRepository) Get(ctx context.Context, id int) (*entity.Business, error) {
	dbBusiness, err := r.db.SQLC.GetBusiness(ctx, int32(id))
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	return convertDBBusinessToEntity(dbBusiness), nil
}

func (r *businessRepository) Update(ctx context.Context, business *entity.Business) error {
	dbBusiness, err := r.db.SQLC.UpdateBusiness(ctx, sqlc.UpdateBusinessParams{
		ID:   int32(business.ID),
		Name: business.Name,
	})
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	business.CreatedAt = dbBusiness.CreatedAt.Time
	return nil
}

func (r *businessRepository) UpdateAppearance(ctx context.Context, id int, logoURL string, colorScheme map[string]interface{}) error {
	colorSchemeJSON, err := json.Marshal(colorScheme)
	if err != nil {
		return fmt.Errorf("failed to marshal color scheme: %w", err)
	}

	_, err = r.db.SQLC.UpdateBusinessAppearance(ctx, sqlc.UpdateBusinessAppearanceParams{
		ID:          int32(id),
		LogoUrl:     pgtype.Text{String: logoURL, Valid: true},
		ColorScheme: colorSchemeJSON,
	})
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	return nil
}

func convertDBBusinessToEntity(dbBusiness sqlc.Business) *entity.Business {
	business := &entity.Business{
		ID:        int(dbBusiness.ID),
		Name:      dbBusiness.Name,
		CreatedAt: dbBusiness.CreatedAt.Time,
	}

	if dbBusiness.LogoUrl.Valid {
		logoURL := dbBusiness.LogoUrl.String
		business.LogoURL = &logoURL
	}

	if len(dbBusiness.ColorScheme) > 0 {
		var colorScheme map[string]interface{}
		if err := json.Unmarshal(dbBusiness.ColorScheme, &colorScheme); err == nil {
			business.ColorScheme = colorScheme
		}
	}

	return business
}
