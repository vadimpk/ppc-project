package repository

type Repositories struct {
	Business    BusinessRepository
	User        UserRepository
	Employee    EmployeeRepository
	Schedule    ScheduleRepository
	Service     BusinessServiceRepository
	Appointment AppointmentRepository
}

func NewRepositories(db *DB) *Repositories {
	return &Repositories{
		Business:    NewBusinessRepository(db),
		User:        NewUserRepository(db),
		Employee:    NewEmployeeRepository(db),
		Schedule:    NewScheduleRepository(db),
		Service:     NewBusinessServiceRepository(db),
		Appointment: NewAppointmentRepository(db),
	}
}
