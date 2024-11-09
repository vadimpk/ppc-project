package entity

import "time"

type Appointment struct {
	ID           int       `json:"id" db:"id"`
	BusinessID   int       `json:"business_id" db:"business_id"`
	ClientID     int       `json:"client_id" db:"client_id"`
	EmployeeID   int       `json:"employee_id" db:"employee_id"`
	ServiceID    int       `json:"service_id" db:"service_id"`
	StartTime    time.Time `json:"start_time" db:"start_time"`
	EndTime      time.Time `json:"end_time" db:"end_time"`
	Status       string    `json:"status" db:"status"`
	ReminderTime *int      `json:"reminder_time" db:"reminder_time"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

const (
	AppointmentStatusScheduled = "scheduled"
	AppointmentStatusCompleted = "completed"
	AppointmentStatusCancelled = "cancelled"
	AppointmentStatusNoShow    = "no_show"
)
