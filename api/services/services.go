package services

import (
	"context"
	"time"

	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

type Services struct {
	Business    BusinessService
	User        UserService
	Employee    EmployeeService
	Schedule    ScheduleService
	Service     BusinessServiceService // renamed to avoid confusion
	Appointment AppointmentService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Business:    NewBusinessService(repos),
		User:        NewUserService(repos),
		Employee:    NewEmployeeService(repos),
		Schedule:    NewScheduleService(repos),
		Service:     NewBusinessServiceService(repos),
		Appointment: NewAppointmentService(repos),
	}
}

// BusinessService handles core business profile management
type BusinessService interface {
	Create(ctx context.Context, business *entity.Business) error
	Get(ctx context.Context, id int) (*entity.Business, error)
	Update(ctx context.Context, business *entity.Business) error
	UpdateAppearance(ctx context.Context, id int, logoURL string, colorScheme map[string]interface{}) error
}

// UserService handles user management and authentication
type UserService interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	CreateBusinessAdmin(ctx context.Context, businessName string, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	// Authentication methods
	Authenticate(ctx context.Context, email string, phone string, password string) (*entity.User, error)
}

// EmployeeService handles employee management
type EmployeeService interface {
	Create(ctx context.Context, employee *entity.Employee) error
	Get(ctx context.Context, id int) (*entity.Employee, error)
	Update(ctx context.Context, employee *entity.Employee) error
	List(ctx context.Context, businessID int) ([]entity.Employee, error)
	AssignServices(ctx context.Context, employeeID int, serviceIDs []int) error
	RemoveServices(ctx context.Context, employeeID int, serviceIDs []int) error
	GetServices(ctx context.Context, employeeID int) ([]entity.BusinessService, error)
}

// BusinessServiceService handles service management
type BusinessServiceService interface {
	Create(ctx context.Context, service *entity.BusinessService) error
	Get(ctx context.Context, id int) (*entity.BusinessService, error)
	Update(ctx context.Context, service *entity.BusinessService) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, businessID int) ([]entity.BusinessService, error)
	ListActive(ctx context.Context, businessID int) ([]entity.BusinessService, error)
}

// ScheduleService handles employee scheduling
type ScheduleService interface {
	// Regular schedule templates
	CreateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error
	UpdateTemplate(ctx context.Context, template *entity.ScheduleTemplate) error
	DeleteTemplate(ctx context.Context, id int) error
	ListTemplates(ctx context.Context, employeeID int) ([]entity.ScheduleTemplate, error)

	// Schedule overrides
	CreateOverride(ctx context.Context, override *entity.ScheduleOverride) error
	UpdateOverride(ctx context.Context, override *entity.ScheduleOverride) error
	DeleteOverride(ctx context.Context, id int) error
	ListOverrides(ctx context.Context, employeeID int, startDate, endDate time.Time) ([]entity.ScheduleOverride, error)

	// Availability checking
	IsAvailable(ctx context.Context, employeeID int, startTime, endTime time.Time) (bool, error)
}

// AppointmentService handles appointment management
type AppointmentService interface {
	Create(ctx context.Context, appointment *entity.Appointment) error
	Get(ctx context.Context, id int) (*entity.Appointment, error)
	Update(ctx context.Context, appointment *entity.Appointment) error
	Cancel(ctx context.Context, id int) error
	ListByBusiness(ctx context.Context, businessID int, startTime, endTime time.Time) ([]entity.Appointment, error)
	ListByEmployee(ctx context.Context, employeeID int, startTime, endTime time.Time) ([]entity.Appointment, error)
	ListByClient(ctx context.Context, clientID int, startTime, endTime time.Time) ([]entity.Appointment, error)
	GetAvailableSlots(ctx context.Context, employeeID int, serviceID int, date time.Time) ([]TimeSlot, error)
}

// Supporting types that match our schema
type TimeSlot struct {
	StartTime time.Time
	EndTime   time.Time
}
