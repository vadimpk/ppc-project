package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository/db/sqlc"
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

func timeToMicroseconds(t time.Time) int64 {
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	duration := t.Sub(midnight)
	return duration.Microseconds()
}

func microsecondsToTime(microseconds int64) time.Time {
	return time.UnixMicro(microseconds).UTC()
}

func (r *scheduleRepository) CreateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error {
	dbTemplate, err := r.db.SQLC.CreateTemplate(ctx, sqlc.CreateTemplateParams{
		EmployeeID: pgtype.Int4{Int32: int32(template.EmployeeID), Valid: true},
		DayOfWeek:  pgtype.Int4{Int32: int32(template.DayOfWeek), Valid: true},
		StartTime:  pgtype.Time{Microseconds: timeToMicroseconds(template.StartTime), Valid: true},
		EndTime:    pgtype.Time{Microseconds: timeToMicroseconds(template.EndTime), Valid: true},
		IsBreak:    pgtype.Bool{Bool: template.IsBreak, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	template.ID = int(dbTemplate.ID)
	template.CreatedAt = dbTemplate.CreatedAt.Time
	return nil
}

func (r *scheduleRepository) UpdateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error {
	dbTemplate, err := r.db.SQLC.UpdateTemplate(ctx, sqlc.UpdateTemplateParams{
		ID:        int32(template.ID),
		StartTime: pgtype.Time{Microseconds: timeToMicroseconds(template.StartTime), Valid: true},
		EndTime:   pgtype.Time{Microseconds: timeToMicroseconds(template.EndTime), Valid: true},
		IsBreak:   pgtype.Bool{Bool: template.IsBreak, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update template: %w", err)
	}

	template.CreatedAt = dbTemplate.CreatedAt.Time
	return nil
}

func (r *scheduleRepository) DeleteTemplate(ctx context.Context, id int) error {
	err := r.db.SQLC.DeleteTemplate(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete template: %w", err)
	}
	return nil
}

func (r *scheduleRepository) ListTemplates(ctx context.Context, employeeID int) ([]entity.ScheduleTemplate, error) {
	dbTemplates, err := r.db.SQLC.ListTemplates(ctx, pgtype.Int4{Int32: int32(employeeID), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	templates := make([]entity.ScheduleTemplate, len(dbTemplates))
	for i, t := range dbTemplates {
		templates[i] = *convertDBTemplateToEntity(t)
	}

	return templates, nil
}

func (r *scheduleRepository) CreateOverride(ctx context.Context, override *entity.ScheduleOverride) error {
	params := sqlc.CreateOverrideParams{
		EmployeeID:   pgtype.Int4{Int32: int32(override.EmployeeID), Valid: true},
		OverrideDate: pgtype.Date{Time: override.OverrideDate, Valid: true},
		IsWorkingDay: pgtype.Bool{Bool: override.IsWorkingDay, Valid: true},
		IsBreak:      pgtype.Bool{Bool: override.IsBreak, Valid: true},
	}

	if override.StartTime != nil {
		params.StartTime = pgtype.Time{
			Microseconds: timeToMicroseconds(*override.StartTime),
			Valid:        true,
		}
	}
	if override.EndTime != nil {
		params.EndTime = pgtype.Time{
			Microseconds: timeToMicroseconds(*override.EndTime),
			Valid:        true,
		}
	}

	dbOverride, err := r.db.SQLC.CreateOverride(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create override: %w", err)
	}

	override.ID = int(dbOverride.ID)
	override.CreatedAt = dbOverride.CreatedAt.Time
	return nil
}

func (r *scheduleRepository) UpdateOverride(ctx context.Context, override *entity.ScheduleOverride) error {
	params := sqlc.UpdateOverrideParams{
		ID:           int32(override.ID),
		IsWorkingDay: pgtype.Bool{Bool: override.IsWorkingDay, Valid: true},
		IsBreak:      pgtype.Bool{Bool: override.IsBreak, Valid: true},
	}

	if override.StartTime != nil {
		params.StartTime = pgtype.Time{
			Microseconds: timeToMicroseconds(*override.StartTime),
			Valid:        true,
		}
	}
	if override.EndTime != nil {
		params.EndTime = pgtype.Time{
			Microseconds: timeToMicroseconds(*override.EndTime),
			Valid:        true,
		}
	}

	dbOverride, err := r.db.SQLC.UpdateOverride(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update override: %w", err)
	}

	override.CreatedAt = dbOverride.CreatedAt.Time
	return nil
}

func (r *scheduleRepository) DeleteOverride(ctx context.Context, id int) error {
	err := r.db.SQLC.DeleteOverride(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete override: %w", err)
	}
	return nil
}

func (r *scheduleRepository) ListOverrides(ctx context.Context, employeeID int, startDate, endDate time.Time) ([]entity.ScheduleOverride, error) {
	params := sqlc.ListOverridesParams{
		EmployeeID:     pgtype.Int4{Int32: int32(employeeID), Valid: true},
		OverrideDate:   pgtype.Date{Time: startDate, Valid: true},
		OverrideDate_2: pgtype.Date{Time: endDate, Valid: true},
	}

	dbOverrides, err := r.db.SQLC.ListOverrides(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list overrides: %w", err)
	}

	overrides := make([]entity.ScheduleOverride, len(dbOverrides))
	for i, o := range dbOverrides {
		overrides[i] = *convertDBOverrideToEntity(o)
	}

	return overrides, nil
}

func convertDBTemplateToEntity(t sqlc.ScheduleTemplate) *entity.ScheduleTemplate {
	return &entity.ScheduleTemplate{
		ID:         int(t.ID),
		EmployeeID: int(t.EmployeeID.Int32),
		DayOfWeek:  int(t.DayOfWeek.Int32),
		StartTime:  microsecondsToTime(t.StartTime.Microseconds),
		EndTime:    microsecondsToTime(t.EndTime.Microseconds),
		IsBreak:    t.IsBreak.Bool,
		CreatedAt:  t.CreatedAt.Time,
	}
}

func convertDBOverrideToEntity(o sqlc.ScheduleOverride) *entity.ScheduleOverride {
	override := &entity.ScheduleOverride{
		ID:           int(o.ID),
		EmployeeID:   int(o.EmployeeID.Int32),
		OverrideDate: o.OverrideDate.Time,
		IsWorkingDay: o.IsWorkingDay.Bool,
		IsBreak:      o.IsBreak.Bool,
		CreatedAt:    o.CreatedAt.Time,
	}

	if o.StartTime.Valid {
		startTime := microsecondsToTime(o.StartTime.Microseconds)
		override.StartTime = &startTime
	}
	if o.EndTime.Valid {
		endTime := microsecondsToTime(o.EndTime.Microseconds)
		override.EndTime = &endTime
	}

	return override
}
