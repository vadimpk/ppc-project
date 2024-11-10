package services

import (
	"context"
	"fmt"
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
	// Validate employee existence and business context
	employee, err := s.repos.Employee.Get(ctx, template.EmployeeID)
	if err != nil {
		return fmt.Errorf("invalid employee: %w", err)
	}
	if !employee.IsActive {
		return fmt.Errorf("employee is not active")
	}

	// Validate template data
	if err := validateTemplateData(template); err != nil {
		return err
	}

	// Check for overlapping templates
	templates, err := s.repos.Schedule.ListTemplates(ctx, template.EmployeeID)
	if err != nil {
		return fmt.Errorf("failed to check existing templates: %w", err)
	}

	for _, t := range templates {
		if t.DayOfWeek == template.DayOfWeek && isTimeOverlap(
			t.StartTime, t.EndTime,
			template.StartTime, template.EndTime,
		) {
			return fmt.Errorf("template overlaps with existing schedule")
		}
	}

	// Create template
	if err := s.repos.Schedule.CreateTemplate(ctx, template); err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}

func (s *scheduleService) UpdateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error {
	// Verify template exists and get original data
	existingTemplates, err := s.repos.Schedule.ListTemplates(ctx, template.EmployeeID)
	if err != nil {
		return fmt.Errorf("failed to get existing templates: %w", err)
	}

	var existing *entity.ScheduleTemplate
	for _, t := range existingTemplates {
		if t.ID == template.ID {
			existing = &t
			break
		}
	}
	if existing == nil {
		return fmt.Errorf("template not found")
	}

	// Validate template data
	if err := validateTemplateData(template); err != nil {
		return err
	}

	// Check for overlapping templates (excluding current template)
	for _, t := range existingTemplates {
		if t.ID != template.ID && t.DayOfWeek == template.DayOfWeek && isTimeOverlap(
			t.StartTime, t.EndTime,
			template.StartTime, template.EndTime,
		) {
			return fmt.Errorf("template overlaps with existing schedule")
		}
	}

	// Update template
	if err := s.repos.Schedule.UpdateTemplate(ctx, template); err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	return nil
}

func (s *scheduleService) DeleteTemplate(ctx context.Context, id int) error {
	if err := s.repos.Schedule.DeleteTemplate(ctx, id); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}
	return nil
}

func (s *scheduleService) ListTemplates(ctx context.Context, employeeID int) ([]entity.ScheduleTemplate, error) {
	// Validate employee existence
	if _, err := s.repos.Employee.Get(ctx, employeeID); err != nil {
		return nil, fmt.Errorf("invalid employee: %w", err)
	}

	templates, err := s.repos.Schedule.ListTemplates(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	return templates, nil
}

func (s *scheduleService) CreateOverride(ctx context.Context, override *entity.ScheduleOverride) error {
	// Validate employee existence and business context
	if _, err := s.repos.Employee.Get(ctx, override.EmployeeID); err != nil {
		return fmt.Errorf("invalid employee: %w", err)
	}

	// Validate override data
	if err := validateOverrideData(override); err != nil {
		return err
	}

	// Check for existing overrides on the same date
	existingOverrides, err := s.repos.Schedule.ListOverrides(ctx, override.EmployeeID, override.OverrideDate, override.OverrideDate)
	if err != nil {
		return fmt.Errorf("failed to check existing overrides: %w", err)
	}

	if override.IsWorkingDay && override.StartTime != nil && override.EndTime != nil {
		for _, o := range existingOverrides {
			if o.StartTime != nil && o.EndTime != nil && isTimeOverlap(
				*o.StartTime, *o.EndTime,
				*override.StartTime, *override.EndTime,
			) {
				return fmt.Errorf("override overlaps with existing schedule")
			}
		}
	}

	// Create override
	if err := s.repos.Schedule.CreateOverride(ctx, override); err != nil {
		return fmt.Errorf("failed to create override: %w", err)
	}

	return nil
}

func (s *scheduleService) UpdateOverride(ctx context.Context, override *entity.ScheduleOverride) error {
	// Validate employee existence and business context
	if _, err := s.repos.Employee.Get(ctx, override.EmployeeID); err != nil {
		return fmt.Errorf("invalid employee: %w", err)
	}

	// Verify override exists and get original data
	existingOverrides, err := s.repos.Schedule.ListOverrides(ctx, override.EmployeeID, override.OverrideDate, override.OverrideDate)
	if err != nil {
		return fmt.Errorf("failed to get existing overrides: %w", err)
	}

	var existing *entity.ScheduleOverride
	for _, o := range existingOverrides {
		if o.ID == override.ID {
			existing = &o
			break
		}
	}
	if existing == nil {
		return fmt.Errorf("override not found")
	}

	// Maintain the original date when updating
	override.OverrideDate = existing.OverrideDate

	// Validate override data
	if err := validateOverrideData(override); err != nil {
		return err
	}

	// Check for overlapping overrides (excluding current override)
	if override.IsWorkingDay && override.StartTime != nil && override.EndTime != nil {
		for _, o := range existingOverrides {
			if o.ID != override.ID && o.IsWorkingDay && o.StartTime != nil && o.EndTime != nil {
				if isTimeOverlap(
					*o.StartTime, *o.EndTime,
					*override.StartTime, *override.EndTime,
				) {
					return fmt.Errorf("override overlaps with existing schedule")
				}
			}
		}
	}

	// Update override
	if err := s.repos.Schedule.UpdateOverride(ctx, override); err != nil {
		return fmt.Errorf("failed to update override: %w", err)
	}

	return nil
}

func (s *scheduleService) DeleteOverride(ctx context.Context, id int) error {
	// Check if override exists first
	overrides, err := s.repos.Schedule.ListOverrides(ctx, 0, time.Now(), time.Now().AddDate(1, 0, 0))
	if err != nil {
		return fmt.Errorf("failed to check override existence: %w", err)
	}

	var exists bool
	for _, o := range overrides {
		if o.ID == id {
			exists = true

			// Check if trying to delete past override
			if o.OverrideDate.Before(time.Now().Truncate(24 * time.Hour)) {
				return fmt.Errorf("cannot delete past overrides")
			}
			break
		}
	}

	if !exists {
		return fmt.Errorf("override not found")
	}

	// Delete override
	if err := s.repos.Schedule.DeleteOverride(ctx, id); err != nil {
		return fmt.Errorf("failed to delete override: %w", err)
	}

	return nil
}

func (s *scheduleService) ListOverrides(ctx context.Context, employeeID int, startDate, endDate time.Time) ([]entity.ScheduleOverride, error) {
	// Validate employee existence
	if _, err := s.repos.Employee.Get(ctx, employeeID); err != nil {
		return nil, fmt.Errorf("invalid employee: %w", err)
	}

	// Validate date range
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Limit the date range to prevent excessive data retrieval
	maxDays := 31
	if endDate.Sub(startDate).Hours()/24 > float64(maxDays) {
		return nil, fmt.Errorf("date range cannot exceed %d days", maxDays)
	}

	// List overrides
	overrides, err := s.repos.Schedule.ListOverrides(ctx, employeeID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list overrides: %w", err)
	}

	return overrides, nil
}

func validateTemplateData(template *entity.ScheduleTemplate) error {
	if template.DayOfWeek < 0 || template.DayOfWeek > 6 {
		return fmt.Errorf("invalid day of week")
	}

	startTimeMinutes := template.StartTime.Hour()*60 + template.StartTime.Minute()
	endTimeMinutes := template.EndTime.Hour()*60 + template.EndTime.Minute()

	if startTimeMinutes >= endTimeMinutes {
		return fmt.Errorf("end time must be after start time")
	}

	return nil
}

func validateOverrideData(override *entity.ScheduleOverride) error {
	if override.OverrideDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return fmt.Errorf("cannot create override for past dates")
	}

	if override.IsWorkingDay {
		if (override.StartTime == nil) != (override.EndTime == nil) {
			return fmt.Errorf("both start time and end time must be provided for working days")
		}
		if override.StartTime != nil && override.EndTime != nil {
			startTimeMinutes := override.StartTime.Hour()*60 + override.StartTime.Minute()
			endTimeMinutes := override.EndTime.Hour()*60 + override.EndTime.Minute()
			if startTimeMinutes >= endTimeMinutes {
				return fmt.Errorf("end time must be after start time")
			}
		}
	}

	return nil
}

func isTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	// Convert to comparable format (minutes since midnight)
	start1Mins := start1.Hour()*60 + start1.Minute()
	end1Mins := end1.Hour()*60 + end1.Minute()
	start2Mins := start2.Hour()*60 + start2.Minute()
	end2Mins := end2.Hour()*60 + end2.Minute()

	return start1Mins < end2Mins && end1Mins > start2Mins
}

func (s *scheduleService) IsAvailable(ctx context.Context, employeeID int, startTime, endTime time.Time) (bool, error) {
	panic("not implemented")
}
