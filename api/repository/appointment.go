package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/vadimpk/ppc-project/entity"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *entity.Appointment) error
	Get(ctx context.Context, id int) (*entity.Appointment, error)
	Update(ctx context.Context, appointment *entity.Appointment) error
	Delete(ctx context.Context, id int) error
	ListByBusiness(ctx context.Context, businessID int, startTime, endTime time.Time) ([]entity.Appointment, error)
	ListByEmployee(ctx context.Context, employeeID int, startTime, endTime time.Time) ([]entity.Appointment, error)
	ListByClient(ctx context.Context, clientID int, startTime, endTime time.Time) ([]entity.Appointment, error)
}

type appointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) AppointmentRepository {
	return &appointmentRepository{
		db: db,
	}
}

func (r *appointmentRepository) Create(ctx context.Context, appointment *entity.Appointment) error {
	panic("not implemented")
}

func (r *appointmentRepository) Get(ctx context.Context, id int) (*entity.Appointment, error) {
	panic("not implemented")
}

func (r *appointmentRepository) Update(ctx context.Context, appointment *entity.Appointment) error {
	panic("not implemented")
}

func (r *appointmentRepository) Delete(ctx context.Context, id int) error {
	panic("not implemented")
}

func (r *appointmentRepository) ListByBusiness(ctx context.Context, businessID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	panic("not implemented")
}

func (r *appointmentRepository) ListByEmployee(ctx context.Context, employeeID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	panic("not implemented")
}

func (r *appointmentRepository) ListByClient(ctx context.Context, clientID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	panic("not implemented")
}
