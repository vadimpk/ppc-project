package services

import (
	"context"
	"fmt"
	"time"

	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

type appointmentService struct {
	repos *repository.Repositories
}

func NewAppointmentService(repos *repository.Repositories) AppointmentService {
	return &appointmentService{
		repos: repos,
	}
}

func (s *appointmentService) Create(ctx context.Context, appointment *entity.Appointment) error {
	// Validate business existence
	if _, err := s.repos.Business.Get(ctx, appointment.BusinessID); err != nil {
		return fmt.Errorf("invalid business: %w", err)
	}

	// Validate client existence
	client, err := s.repos.User.Get(ctx, appointment.ClientID)
	if err != nil {
		return fmt.Errorf("invalid client: %w", err)
	}
	if client.Role != entity.RoleClient {
		return fmt.Errorf("user is not a client")
	}

	// Validate employee existence and active status
	employee, err := s.repos.Employee.Get(ctx, appointment.EmployeeID)
	if err != nil {
		return fmt.Errorf("invalid employee: %w", err)
	}
	if !employee.IsActive {
		return fmt.Errorf("employee is not active")
	}

	// Validate service existence and availability for employee
	service, err := s.repos.Service.Get(ctx, appointment.ServiceID)
	if err != nil {
		return fmt.Errorf("invalid service: %w", err)
	}
	if !service.IsActive {
		return fmt.Errorf("service is not active")
	}

	// Validate service assignment to employee
	employeeServices, err := s.repos.Employee.GetServices(ctx, appointment.EmployeeID)
	if err != nil {
		return fmt.Errorf("failed to check employee services: %w", err)
	}

	var serviceAssigned bool
	for _, s := range employeeServices {
		if s.ID == appointment.ServiceID {
			serviceAssigned = true
			break
		}
	}
	if !serviceAssigned {
		return fmt.Errorf("service is not assigned to employee")
	}

	appointment.EndTime = appointment.StartTime.Add(time.Duration(service.Duration) * time.Minute)
	// Validate appointment time
	if err := s.validateAppointmentTime(ctx, appointment, service.Duration); err != nil {
		return fmt.Errorf("invalid appointment time: %w", err)
	}

	// Set default status
	appointment.Status = entity.AppointmentStatusScheduled

	// Create appointment
	if err := s.repos.Appointment.Create(ctx, appointment); err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}

	return nil
}

func (s *appointmentService) Get(ctx context.Context, id int) (*entity.Appointment, error) {
	appointment, err := s.repos.Appointment.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return appointment, nil
}

func (s *appointmentService) Update(ctx context.Context, appointment *entity.Appointment) error {
	panic("implement me")
	//// Verify appointment exists and get current data
	//existing, err := s.repos.Appointment.Get(ctx, appointment.ID)
	//if err != nil {
	//	return fmt.Errorf("invalid appointment: %w", err)
	//}
	//
	//// Only allow updates for scheduled appointments
	//if existing.Status != entity.AppointmentStatusScheduled {
	//	return fmt.Errorf("can only update scheduled appointments")
	//}
	//
	//// Cannot update past appointments
	//if existing.StartTime.Before(time.Now()) {
	//	return fmt.Errorf("cannot update past appointments")
	//}
	//
	//// Maintain original IDs and creation data
	//appointment.BusinessID = existing.BusinessID
	//appointment.ClientID = existing.ClientID
	//appointment.EmployeeID = existing.EmployeeID
	//appointment.ServiceID = existing.ServiceID
	//appointment.CreatedAt = existing.CreatedAt
	//
	//// Get service duration for validation
	//service, err := s.repos.Service.Get(ctx, appointment.ServiceID)
	//if err != nil {
	//	return fmt.Errorf("failed to get service details: %w", err)
	//}
	//
	//// Validate new appointment time
	//if err := s.validateAppointmentTime(ctx, appointment, service.Duration); err != nil {
	//	return fmt.Errorf("invalid appointment time: %w", err)
	//}
	//
	//// Update appointment
	//if err := s.repos.Appointment.Update(ctx, appointment); err != nil {
	//	return fmt.Errorf("failed to update appointment: %w", err)
	//}
	//
	//return nil
}

func (s *appointmentService) Cancel(ctx context.Context, id int) error {
	// Verify appointment exists
	appointment, err := s.repos.Appointment.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("invalid appointment: %w", err)
	}

	// Only allow cancellation of scheduled appointments
	if appointment.Status != entity.AppointmentStatusScheduled {
		return fmt.Errorf("can only cancel scheduled appointments")
	}

	// Cannot cancel past appointments
	if appointment.StartTime.Before(time.Now()) {
		return fmt.Errorf("cannot cancel past appointments")
	}

	// Cancel appointment
	if err := s.repos.Appointment.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}

	return nil
}

func (s *appointmentService) ListByBusiness(ctx context.Context, businessID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	// Validate business existence
	if _, err := s.repos.Business.Get(ctx, businessID); err != nil {
		return nil, fmt.Errorf("invalid business: %w", err)
	}

	if err := validateDateRange(startTime, endTime); err != nil {
		return nil, err
	}

	appointments, err := s.repos.Appointment.ListByBusiness(ctx, businessID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}

	return appointments, nil
}

func (s *appointmentService) ListByClient(ctx context.Context, clientID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	// Validate client existence
	if _, err := s.repos.User.Get(ctx, clientID); err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}

	if err := validateDateRange(startTime, endTime); err != nil {
		return nil, err
	}

	appointments, err := s.repos.Appointment.ListByClient(ctx, clientID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}

	return appointments, nil
}

func (s *appointmentService) ListByEmployee(ctx context.Context, employeeID int, startTime time.Time, endTime time.Time) ([]entity.Appointment, error) {
	// Validate client existence
	if _, err := s.repos.Employee.Get(ctx, employeeID); err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}

	if err := validateDateRange(startTime, endTime); err != nil {
		return nil, err
	}

	appointments, err := s.repos.Appointment.ListByEmployee(ctx, employeeID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}

	return appointments, nil
}

func (s *appointmentService) GetAvailableSlots(ctx context.Context, employeeID int, serviceID int, date time.Time) ([]TimeSlot, error) {
	// Validate employee and service
	employee, err := s.repos.Employee.Get(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee: %w", err)
	}
	if !employee.IsActive {
		return nil, fmt.Errorf("employee is not active")
	}

	service, err := s.repos.Service.Get(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("invalid service: %w", err)
	}
	if !service.IsActive {
		return nil, fmt.Errorf("service is not active")
	}

	// Check if service is assigned to employee
	employeeServices, err := s.repos.Employee.GetServices(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check employee services: %w", err)
	}
	var serviceAssigned bool
	for _, s := range employeeServices {
		if s.ID == serviceID {
			serviceAssigned = true
			break
		}
	}
	if !serviceAssigned {
		return nil, fmt.Errorf("service is not assigned to employee")
	}

	// Get schedule for the date
	schedules, err := s.repos.Schedule.GetEmployeeSchedule(ctx, employeeID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee schedule: %w", err)
	}

	if schedules == nil {
		return nil, fmt.Errorf("no schedule found for the date")
	}

	// Get existing appointments
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	appointments, err := s.repos.Appointment.ListByEmployee(ctx, employeeID, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing appointments: %w", err)
	}

	// Generate available slots
	slots := generateAvailableSlots(schedules, appointments, service.Duration, date)
	return slots, nil
}

func (s *appointmentService) validateAppointmentTime(ctx context.Context, appointment *entity.Appointment, serviceDuration int) error {
	// Appointment must be in the future
	if appointment.StartTime.Before(time.Now()) {
		return fmt.Errorf("cannot create appointments in the past")
	}

	// Calculate end time based on service duration if not provided
	if appointment.EndTime.IsZero() {
		appointment.EndTime = appointment.StartTime.Add(time.Duration(serviceDuration) * time.Minute)
	}

	// Validate time slot duration matches service duration
	duration := int(appointment.EndTime.Sub(appointment.StartTime).Minutes())
	if duration != serviceDuration {
		return fmt.Errorf("appointment duration must match service duration")
	}

	// Check employee schedule
	date := appointment.StartTime.Truncate(24 * time.Hour)
	schedules, err := s.repos.Schedule.GetEmployeeSchedule(ctx, appointment.EmployeeID, date)
	if err != nil {
		return fmt.Errorf("failed to check employee schedule: %w", err)
	}

	if schedules == nil || !isTimeSlotInSchedule(appointment.StartTime, appointment.EndTime, schedules) {
		return fmt.Errorf("appointment time is outside employee's working hours")
	}

	// Check for overlapping appointments
	isAvailable, err := s.repos.Appointment.IsEmployeeAvailable(ctx, appointment.EmployeeID, appointment.StartTime, appointment.EndTime)
	if err != nil {
		return fmt.Errorf("failed to check employee availability: %w", err)
	}
	if !isAvailable {
		return fmt.Errorf("time slot is not available")
	}

	return nil
}

func validateDateRange(startTime, endTime time.Time) error {
	if endTime.Before(startTime) {
		return fmt.Errorf("end time must be after start time")
	}

	maxRange := 31 * 24 * time.Hour // 31 days
	if endTime.Sub(startTime) > maxRange {
		return fmt.Errorf("date range cannot exceed 31 days")
	}

	return nil
}

func isTimeSlotInSchedule(start, end time.Time, schedule *entity.ScheduleTemplate) bool {
	startTime := time.Date(start.Year(), start.Month(), start.Day(), schedule.StartTime.Hour(), schedule.StartTime.Minute(), 0, 0, start.Location()).Add(-1 * time.Minute)
	endTime := time.Date(end.Year(), end.Month(), end.Day(), schedule.EndTime.Hour(), schedule.EndTime.Minute(), 0, 0, end.Location()).Add(1 * time.Minute)
	if !schedule.IsBreak && start.After(startTime) && end.Before(endTime) {
		return true
	}
	return false
}

func generateAvailableSlots(schedule *entity.ScheduleTemplate, appointments []entity.Appointment, duration int, date time.Time) []TimeSlot {
	if schedule.IsBreak {
		return nil
	}

	var availableSlots []TimeSlot
	slotDuration := time.Duration(duration) * time.Minute

	now := time.Now().Add(15 * time.Minute).UTC()

	// Start iterating from the start time of the schedule to the end time
	currentStartTime := time.Date(date.Year(), date.Month(), date.Day(), schedule.StartTime.Hour(), schedule.StartTime.Minute(), 0, 0, now.Location())
	endTime := time.Date(date.Year(), date.Month(), date.Day(), schedule.EndTime.Hour(), schedule.EndTime.Minute(), 0, 0, now.Location())
	for currentStartTime.Add(slotDuration).Before(endTime) || currentStartTime.Add(slotDuration).Equal(endTime) {
		currentEndTime := currentStartTime.Add(slotDuration)
		if currentEndTime.Before(now) {
			continue
		}

		// Check if this slot conflicts with any appointment
		conflict := false
		for _, appointment := range appointments {
			if appointment.StartTime.Before(currentEndTime) && appointment.EndTime.After(currentStartTime) {
				conflict = true
				break
			}
		}

		// If no conflict, add the slot to available slots
		if !conflict {
			availableSlots = append(availableSlots, TimeSlot{
				StartTime: currentStartTime,
				EndTime:   currentEndTime,
			})
		}

		// Move to the next time slot
		currentStartTime = currentEndTime
	}

	return availableSlots
}
