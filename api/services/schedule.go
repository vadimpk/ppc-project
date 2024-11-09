package services

import (
	"context"
	"time"

	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

type scheduleService struct {
	repos *repository.Repositories
}

func NewScheduleService(repos *repository.Repositories) ScheduleService {
	return &scheduleService{
		repos: repos,
	}
}

func (s *scheduleService) CreateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error {
	panic("not implemented")
}

func (s *scheduleService) UpdateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error {
	panic("not implemented")
}

func (s *scheduleService) DeleteTemplate(ctx context.Context, id int) error {
	panic("not implemented")
}

func (s *scheduleService) ListTemplates(ctx context.Context, employeeID int) ([]entity.ScheduleTemplate, error) {
	panic("not implemented")
}

func (s *scheduleService) CreateOverride(ctx context.Context, override *entity.ScheduleOverride) error {
	panic("not implemented")
}

func (s *scheduleService) UpdateOverride(ctx context.Context, override *entity.ScheduleOverride) error {
	panic("not implemented")
}

func (s *scheduleService) DeleteOverride(ctx context.Context, id int) error {
	panic("not implemented")
}

func (s *scheduleService) ListOverrides(ctx context.Context, employeeID int, startDate, endDate time.Time) ([]entity.ScheduleOverride, error) {
	panic("not implemented")
}

func (s *scheduleService) IsAvailable(ctx context.Context, employeeID int, startTime, endTime time.Time) (bool, error) {
	panic("not implemented")
}
