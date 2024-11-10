package controller

import (
	"github.com/vadimpk/ppc-project/pkg/auth"
	"github.com/vadimpk/ppc-project/services"
)

type Handlers struct {
	Business    *BusinessHandler
	User        *UserHandler
	Employee    *EmployeeHandler
	Service     *BusinessServiceHandler
	Schedule    *ScheduleHandler
	Appointment *AppointmentHandler
}

func NewHandlers(services *services.Services, tokenManager *auth.TokenManager) *Handlers {
	return &Handlers{
		Business:    NewBusinessHandler(services.Business),
		User:        NewUserHandler(services.User, tokenManager),
		Employee:    NewEmployeeHandler(services.Employee),
		Service:     NewBusinessServiceHandler(services.Service),
		Schedule:    NewScheduleHandler(services.Schedule),
		Appointment: NewAppointmentHandler(services.Appointment),
	}
}
