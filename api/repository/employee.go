package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository/db/sqlc"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --dir . --name EmployeeRepository --output ./mocks
type EmployeeRepository interface {
	Create(ctx context.Context, employee *entity.Employee) error
	Get(ctx context.Context, id int) (*entity.Employee, error)
	Update(ctx context.Context, employee *entity.Employee) error
	List(ctx context.Context, businessID int) ([]entity.Employee, error)
	AssignServices(ctx context.Context, employeeID int, serviceIDs []int) error
	RemoveServices(ctx context.Context, employeeID int, serviceIDs []int) error
	GetServices(ctx context.Context, employeeID int) ([]entity.BusinessService, error)
	GetIDByUserID(ctx context.Context, userID int) (int, error)
	ListByServiceID(ctx context.Context, serviceID int) ([]entity.Employee, error)
}

type employeeRepository struct {
	db *DB
}

func NewEmployeeRepository(db *DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (r *employeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
	var specialization pgtype.Text
	if employee.Specialization != nil {
		specialization = r.db.ValidText(*employee.Specialization)
	}

	dbEmployee, err := r.db.SQLC.CreateEmployee(ctx, sqlc.CreateEmployeeParams{
		BusinessID:     pgtype.Int4{Int32: int32(employee.BusinessID), Valid: true},
		UserID:         pgtype.Int4{Int32: int32(employee.UserID), Valid: true},
		Specialization: specialization,
		IsActive:       pgtype.Bool{Bool: employee.IsActive, Valid: true},
	})
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	employee.ID = int(dbEmployee.ID)
	employee.CreatedAt = dbEmployee.CreatedAt.Time
	return nil
}

func (r *employeeRepository) Get(ctx context.Context, id int) (*entity.Employee, error) {
	dbEmployee, err := r.db.SQLC.GetEmployee(ctx, int32(id))
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	return convertDBEmployeeToEntity(dbEmployee), nil
}

func (r *employeeRepository) GetIDByUserID(ctx context.Context, userID int) (int, error) {
	dbEmployee, err := r.db.SQLC.GetEmployeeIDByUserID(ctx, pgtype.Int4{Int32: int32(userID), Valid: true})
	if err != nil {
		return 0, r.db.HandleBasicErrors(err)
	}

	return int(dbEmployee), nil
}

func (r *employeeRepository) Update(ctx context.Context, employee *entity.Employee) error {
	var specialization pgtype.Text
	if employee.Specialization != nil {
		specialization = r.db.ValidText(*employee.Specialization)
	}

	dbEmployee, err := r.db.SQLC.UpdateEmployee(ctx, sqlc.UpdateEmployeeParams{
		ID:             int32(employee.ID),
		Specialization: specialization,
		IsActive:       pgtype.Bool{Bool: employee.IsActive, Valid: true},
	})
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	employee.CreatedAt = dbEmployee.CreatedAt.Time
	return nil
}

func (r *employeeRepository) List(ctx context.Context, businessID int) ([]entity.Employee, error) {
	dbEmployees, err := r.db.SQLC.ListEmployees(ctx, pgtype.Int4{Int32: int32(businessID), Valid: true})
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	employees := make([]entity.Employee, len(dbEmployees))
	for i, dbEmployee := range dbEmployees {
		employees[i] = *convertDBEmployeeToEntity(sqlc.GetEmployeeRow(dbEmployee))
	}

	return employees, nil
}

func (r *employeeRepository) ListByServiceID(ctx context.Context, serviceID int) ([]entity.Employee, error) {
	dbEmployees, err := r.db.SQLC.ListEmployeesByService(ctx, int32(serviceID))
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	employees := make([]entity.Employee, len(dbEmployees))
	for i, dbEmployee := range dbEmployees {
		employees[i] = *convertDBEmployeeToEntity(sqlc.GetEmployeeRow(dbEmployee))
	}

	return employees, nil
}

func (r *employeeRepository) AssignServices(ctx context.Context, employeeID int, serviceIDs []int) error {
	// Convert []int to []int32
	serviceIDsInt32 := make([]int32, len(serviceIDs))
	for i, id := range serviceIDs {
		serviceIDsInt32[i] = int32(id)
	}

	err := r.db.SQLC.AssignServices(ctx, sqlc.AssignServicesParams{
		EmployeeID: int32(employeeID),
		Column2:    serviceIDsInt32,
	})
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	return nil
}

func (r *employeeRepository) RemoveServices(ctx context.Context, employeeID int, serviceIDs []int) error {
	// Convert []int to []int32
	serviceIDsInt32 := make([]int32, len(serviceIDs))
	for i, id := range serviceIDs {
		serviceIDsInt32[i] = int32(id)
	}

	err := r.db.SQLC.RemoveServices(ctx, sqlc.RemoveServicesParams{
		EmployeeID: int32(employeeID),
		Column2:    serviceIDsInt32,
	})
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	return nil
}

func (r *employeeRepository) GetServices(ctx context.Context, employeeID int) ([]entity.BusinessService, error) {
	dbServices, err := r.db.SQLC.GetEmployeeServices(ctx, int32(employeeID))
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	services := make([]entity.BusinessService, len(dbServices))
	for i, dbService := range dbServices {
		services[i] = *convertDBServiceToEntity(dbService)
	}

	return services, nil
}

func convertDBEmployeeToEntity(row sqlc.GetEmployeeRow) *entity.Employee {
	employee := &entity.Employee{
		ID:         int(row.ID),
		BusinessID: int(row.BusinessID.Int32),
		UserID:     int(row.UserID.Int32),
		IsActive:   row.IsActive.Bool,
		CreatedAt:  row.CreatedAt.Time,
		User: &entity.User{
			ID:         int(row.UserID.Int32),
			BusinessID: int(row.BusinessID.Int32),
			FullName:   row.UserFullName,
			Role:       row.UserRole,
			CreatedAt:  row.UserCreatedAt.Time,
		},
	}

	if row.Specialization.Valid {
		spec := row.Specialization.String
		employee.Specialization = &spec
	}

	if row.UserEmail.Valid {
		email := row.UserEmail.String
		employee.User.Email = &email
	}

	if row.UserPhone.Valid {
		phone := row.UserPhone.String
		employee.User.Phone = &phone
	}

	return employee
}
