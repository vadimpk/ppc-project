package entity

import "time"

type ScheduleTemplate struct {
	ID         int       `json:"id" db:"id"`
	EmployeeID int       `json:"employee_id" db:"employee_id"`
	DayOfWeek  int       `json:"day_of_week" db:"day_of_week"`
	StartTime  time.Time `json:"start_time" db:"start_time"`
	EndTime    time.Time `json:"end_time" db:"end_time"`
	IsBreak    bool      `json:"is_break" db:"is_break"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type ScheduleOverride struct {
	ID           int        `json:"id" db:"id"`
	EmployeeID   int        `json:"employee_id" db:"employee_id"`
	OverrideDate time.Time  `json:"override_date" db:"override_date"`
	StartTime    *time.Time `json:"start_time" db:"start_time"`
	EndTime      *time.Time `json:"end_time" db:"end_time"`
	IsWorkingDay bool       `json:"is_working_day" db:"is_working_day"`
	IsBreak      bool       `json:"is_break" db:"is_break"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}
