package services

import (
	"context"
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
	panic("not implemented")
}

func (s *appointmentService) Get(ctx context.Context, id int) (*entity.Appointment, error) {
	panic("not implemented")
}

func (s *appointmentService) Update(ctx context.Context, appointment *entity.Appointment) error {
	panic("not implemented")
}

func (s *appointmentService) Cancel(ctx context.Context, id int) error {
	panic("not implemented")
}

func (s *appointmentService) ListByBusiness(ctx context.Context, businessID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	panic("not implemented")
}

func (s *appointmentService) ListByEmployee(ctx context.Context, employeeID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	panic("not implemented")
}

func (s *appointmentService) ListByClient(ctx context.Context, clientID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	panic("not implemented")
}

func (s *appointmentService) GetAvailableSlots(ctx context.Context, employeeID int, serviceID int, date time.Time) ([]TimeSlot, error) {
	panic("not implemented")
}
