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

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *entity.Appointment) error
	Get(ctx context.Context, id int) (*entity.Appointment, error)
	Update(ctx context.Context, appointment *entity.Appointment) error
	Delete(ctx context.Context, id int) error
	ListByBusiness(ctx context.Context, businessID int, startTime, endTime time.Time) ([]entity.Appointment, error)
	ListByEmployee(ctx context.Context, employeeID int, startTime, endTime time.Time) ([]entity.Appointment, error)
	ListByClient(ctx context.Context, clientID int, startTime, endTime time.Time) ([]entity.Appointment, error)
	IsEmployeeAvailable(ctx context.Context, employeeID int, startTime, endTime time.Time) (bool, error)
}

type appointmentRepository struct {
	db *DB
}

func NewAppointmentRepository(db *DB) AppointmentRepository {
	return &appointmentRepository{
		db: db,
	}
}

func (r *appointmentRepository) Create(ctx context.Context, appointment *entity.Appointment) error {
	var reminderTime pgtype.Int4
	if appointment.ReminderTime != nil {
		reminderTime = pgtype.Int4{Int32: int32(*appointment.ReminderTime), Valid: true}
	}

	dbAppointment, err := r.db.SQLC.CreateAppointment(ctx, sqlc.CreateAppointmentParams{
		BusinessID:   pgtype.Int4{Int32: int32(appointment.BusinessID), Valid: true},
		ClientID:     pgtype.Int4{Int32: int32(appointment.ClientID), Valid: true},
		EmployeeID:   pgtype.Int4{Int32: int32(appointment.EmployeeID), Valid: true},
		ServiceID:    pgtype.Int4{Int32: int32(appointment.ServiceID), Valid: true},
		StartTime:    pgtype.Timestamptz{Time: appointment.StartTime, Valid: true},
		EndTime:      pgtype.Timestamptz{Time: appointment.EndTime, Valid: true},
		Status:       pgtype.Text{String: appointment.Status, Valid: true},
		ReminderTime: reminderTime,
	})
	if err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}

	appointment.ID = int(dbAppointment.ID)
	appointment.CreatedAt = dbAppointment.CreatedAt.Time
	return nil
}

func (r *appointmentRepository) Get(ctx context.Context, id int) (*entity.Appointment, error) {
	dbAppointment, err := r.db.SQLC.GetAppointment(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return convertDBAppointmentToEntity(dbAppointment), nil
}

func (r *appointmentRepository) Update(ctx context.Context, appointment *entity.Appointment) error {
	var reminderTime pgtype.Int4
	if appointment.ReminderTime != nil {
		reminderTime = pgtype.Int4{Int32: int32(*appointment.ReminderTime), Valid: true}
	}

	dbAppointment, err := r.db.SQLC.UpdateAppointment(ctx, sqlc.UpdateAppointmentParams{
		ID:           int32(appointment.ID),
		StartTime:    pgtype.Timestamptz{Time: appointment.StartTime, Valid: true},
		EndTime:      pgtype.Timestamptz{Time: appointment.EndTime, Valid: true},
		Status:       pgtype.Text{String: appointment.Status, Valid: true},
		ReminderTime: reminderTime,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	appointment.CreatedAt = dbAppointment.CreatedAt.Time
	return nil
}

func (r *appointmentRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.SQLC.CancelAppointment(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}
	return nil
}

func (r *appointmentRepository) ListByBusiness(ctx context.Context, businessID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	dbAppointments, err := r.db.SQLC.ListBusinessAppointments(ctx, sqlc.ListBusinessAppointmentsParams{
		BusinessID:  pgtype.Int4{Int32: int32(businessID), Valid: true},
		StartTime:   pgtype.Timestamptz{Time: startTime, Valid: true},
		StartTime_2: pgtype.Timestamptz{Time: endTime, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list business appointments: %w", err)
	}

	appointments := make([]entity.Appointment, len(dbAppointments))
	for i, a := range dbAppointments {
		appointments[i] = *convertDBAppointmentToEntity(sqlc.GetAppointmentRow(a))
	}

	return appointments, nil
}

func (r *appointmentRepository) ListByEmployee(ctx context.Context, employeeID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	dbAppointments, err := r.db.SQLC.ListEmployeeAppointments(ctx, sqlc.ListEmployeeAppointmentsParams{
		EmployeeID:  pgtype.Int4{Int32: int32(employeeID), Valid: true},
		StartTime:   pgtype.Timestamptz{Time: startTime, Valid: true},
		StartTime_2: pgtype.Timestamptz{Time: endTime, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list employee appointments: %w", err)
	}

	appointments := make([]entity.Appointment, len(dbAppointments))
	for i, a := range dbAppointments {
		appointments[i] = *convertDBAppointmentToEntity(sqlc.GetAppointmentRow(a))
	}

	return appointments, nil
}

func (r *appointmentRepository) ListByClient(ctx context.Context, clientID int, startTime, endTime time.Time) ([]entity.Appointment, error) {
	dbAppointments, err := r.db.SQLC.ListClientAppointments(ctx, sqlc.ListClientAppointmentsParams{
		ClientID:    pgtype.Int4{Int32: int32(clientID), Valid: true},
		StartTime:   pgtype.Timestamptz{Time: startTime, Valid: true},
		StartTime_2: pgtype.Timestamptz{Time: endTime, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list client appointments: %w", err)
	}

	appointments := make([]entity.Appointment, len(dbAppointments))
	for i, a := range dbAppointments {
		appointments[i] = *convertDBAppointmentToEntity(sqlc.GetAppointmentRow(a))
	}

	return appointments, nil
}

func (r *appointmentRepository) IsEmployeeAvailable(ctx context.Context, employeeID int, startTime, endTime time.Time) (bool, error) {
	available, err := r.db.SQLC.CheckEmployeeAvailability(ctx, sqlc.CheckEmployeeAvailabilityParams{
		EmployeeID: pgtype.Int4{Int32: int32(employeeID), Valid: true},
		Overlaps:   pgtype.Timestamptz{Time: startTime, Valid: true},
		Overlaps_2: pgtype.Timestamptz{Time: endTime, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("failed to check employee availability: %w", err)
	}

	return available, nil
}

type AppointmentWithDetails struct {
	entity.Appointment
	Client   *entity.User            `json:"client"`
	Employee *entity.User            `json:"employee"`
	Service  *entity.BusinessService `json:"service"`
}

func convertDBAppointmentToEntity(a sqlc.GetAppointmentRow) *entity.Appointment {
	appointment := &entity.Appointment{
		ID:         int(a.ID),
		BusinessID: int(a.BusinessID.Int32),
		ClientID:   int(a.ClientID.Int32),
		EmployeeID: int(a.EmployeeID.Int32),
		ServiceID:  int(a.ServiceID.Int32),
		StartTime:  a.StartTime.Time,
		EndTime:    a.EndTime.Time,
		Status:     a.Status.String,
		CreatedAt:  a.CreatedAt.Time,
	}

	if a.ReminderTime.Valid {
		reminderTime := int(a.ReminderTime.Int32)
		appointment.ReminderTime = &reminderTime
	}

	// Add client details
	appointment.Client = &entity.User{
		FullName: a.ClientFullName,
	}
	if a.ClientEmail.Valid {
		email := a.ClientEmail.String
		appointment.Client.Email = &email
	}
	if a.ClientPhone.Valid {
		phone := a.ClientPhone.String
		appointment.Client.Phone = &phone
	}

	// Add employee details
	appointment.Employee = &entity.User{
		FullName: a.EmployeeFullName,
	}
	if a.EmployeeEmail.Valid {
		email := a.EmployeeEmail.String
		appointment.Employee.Email = &email
	}
	if a.EmployeePhone.Valid {
		phone := a.EmployeePhone.String
		appointment.Employee.Phone = &phone
	}

	// Add service details
	appointment.Service = &entity.BusinessService{
		Name:     a.ServiceName,
		Duration: int(a.ServiceDuration),
		Price:    int(a.ServicePrice),
	}

	return appointment
}
