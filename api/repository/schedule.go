package repository

import (
	"context"
	"time"

	"github.com/vadimpk/ppc-project/entity"
)

type ScheduleRepository interface {
	CreateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error
	UpdateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error
	DeleteTemplate(ctx context.Context, id int) error
	ListTemplates(ctx context.Context, employeeID int) ([]entity.ScheduleTemplate, error)
	CreateOverride(ctx context.Context, override *entity.ScheduleOverride) error
	UpdateOverride(ctx context.Context, override *entity.ScheduleOverride) error
	DeleteOverride(ctx context.Context, id int) error
	ListOverrides(ctx context.Context, employeeID int, startDate, endDate time.Time) ([]entity.ScheduleOverride, error)
}

type scheduleRepository struct {
	db *DB
}

func NewScheduleRepository(db *DB) ScheduleRepository {
	return &scheduleRepository{
		db: db,
	}
}

func (r *scheduleRepository) CreateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error {
	panic("not implemented")
}

func (r *scheduleRepository) UpdateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error {
	panic("not implemented")
}

func (r *scheduleRepository) DeleteTemplate(ctx context.Context, id int) error {
	panic("not implemented")
}

func (r *scheduleRepository) ListTemplates(ctx context.Context, employeeID int) ([]entity.ScheduleTemplate, error) {
	panic("not implemented")
}

func (r *scheduleRepository) CreateOverride(ctx context.Context, override *entity.ScheduleOverride) error {
	panic("not implemented")
}

func (r *scheduleRepository) UpdateOverride(ctx context.Context, override *entity.ScheduleOverride) error {
	panic("not implemented")
}

func (r *scheduleRepository) DeleteOverride(ctx context.Context, id int) error {
	panic("not implemented")
}

func (r *scheduleRepository) ListOverrides(ctx context.Context, employeeID int, startDate, endDate time.Time) ([]entity.ScheduleOverride, error) {
	panic("not implemented")
}
