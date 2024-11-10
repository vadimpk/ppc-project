package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
	"github.com/vadimpk/ppc-project/repository/mocks"
	"github.com/vadimpk/ppc-project/services"
)

func TestEmployeeService_Create(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
		employeeRepo *mocks.EmployeeRepository
		userRepo     *mocks.UserRepository
	}

	type args struct {
		employee *entity.Employee
	}

	type expected struct {
		err error
	}

	businessID := 1
	userID := 1

	employee := &entity.Employee{
		BusinessID:     businessID,
		UserID:         userID,
		Specialization: stringPtr("Hair Stylist"),
		IsActive:       true,
	}

	user := &entity.User{
		ID:         userID,
		BusinessID: businessID,
		Role:       entity.RoleClient,
	}

	userAsEmployee := &entity.User{
		ID:         userID,
		BusinessID: businessID,
		Role:       entity.RoleEmployee,
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: employee successfully created",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.userRepo.On("Get", ctx, userID).Return(user, nil)
				m.userRepo.On("Update", ctx, userAsEmployee).Return(nil)
				m.employeeRepo.On("Create", ctx, employee).Return(nil)
			},
			args: args{
				employee: employee,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "positive: employee created with existing employee role",
			mock: func(m mocksForExecution) {
				existingEmployee := user
				existingEmployee.Role = entity.RoleEmployee
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.userRepo.On("Get", ctx, userID).Return(existingEmployee, nil)
				m.employeeRepo.On("Create", ctx, employee).Return(nil)
			},
			args: args{
				employee: employee,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "negative: business not found",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(nil, repository.ErrNotFound)
			},
			args: args{
				employee: employee,
			},
			expected: expected{
				err: fmt.Errorf("invalid business: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: user not found",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.userRepo.On("Get", ctx, userID).Return(nil, repository.ErrNotFound)
			},
			args: args{
				employee: employee,
			},
			expected: expected{
				err: fmt.Errorf("invalid user: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: user from different business",
			mock: func(m mocksForExecution) {
				userFromDifferentBusiness := user
				userFromDifferentBusiness.BusinessID = 999
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.userRepo.On("Get", ctx, userID).Return(userFromDifferentBusiness, nil)
			},
			args: args{
				employee: employee,
			},
			expected: expected{
				err: fmt.Errorf("user does not belong to the business"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			// Init mocks
			businessRepoMock := mocks.NewBusinessRepository(t)
			employeeRepoMock := mocks.NewEmployeeRepository(t)
			userRepoMock := mocks.NewUserRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
				employeeRepo: employeeRepoMock,
				userRepo:     userRepoMock,
			})

			// Init service
			employeeService := services.NewEmployeeService(&repository.Repositories{
				Business: businessRepoMock,
				Employee: employeeRepoMock,
				User:     userRepoMock,
			})

			// Execute
			err := employeeService.Create(ctx, tc.args.employee)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeService_AssignServices(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		employeeRepo *mocks.EmployeeRepository
		serviceRepo  *mocks.BusinessServiceRepository
	}

	type args struct {
		employeeID int
		serviceIDs []int
	}

	type expected struct {
		err error
	}

	employeeID := 1
	businessID := 1
	serviceIDs := []int{1, 2}

	employee := &entity.Employee{
		ID:         employeeID,
		BusinessID: businessID,
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: services successfully assigned",
			mock: func(m mocksForExecution) {
				m.employeeRepo.On("Get", ctx, employeeID).Return(employee, nil)
				for _, serviceID := range serviceIDs {
					m.serviceRepo.On("Get", ctx, serviceID).Return(&entity.BusinessService{
						ID:         serviceID,
						BusinessID: businessID,
					}, nil)
				}
				m.employeeRepo.On("AssignServices", ctx, employeeID, serviceIDs).Return(nil)
			},
			args: args{
				employeeID: employeeID,
				serviceIDs: serviceIDs,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "negative: employee not found",
			mock: func(m mocksForExecution) {
				m.employeeRepo.On("Get", ctx, employeeID).Return(nil, repository.ErrNotFound)
			},
			args: args{
				employeeID: employeeID,
				serviceIDs: serviceIDs,
			},
			expected: expected{
				err: fmt.Errorf("invalid employee: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: service from different business",
			mock: func(m mocksForExecution) {
				m.employeeRepo.On("Get", ctx, employeeID).Return(employee, nil)
				m.serviceRepo.On("Get", ctx, serviceIDs[0]).Return(&entity.BusinessService{
					ID:         serviceIDs[0],
					BusinessID: 999, // different business
				}, nil)
			},
			args: args{
				employeeID: employeeID,
				serviceIDs: serviceIDs,
			},
			expected: expected{
				err: fmt.Errorf("service %d does not belong to the employee's business", serviceIDs[0]),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			employeeRepoMock := mocks.NewEmployeeRepository(t)
			serviceRepoMock := mocks.NewBusinessServiceRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				employeeRepo: employeeRepoMock,
				serviceRepo:  serviceRepoMock,
			})

			// Init service
			employeeService := services.NewEmployeeService(&repository.Repositories{
				Employee: employeeRepoMock,
				Service:  serviceRepoMock,
			})

			// Execute
			err := employeeService.AssignServices(ctx, tc.args.employeeID, tc.args.serviceIDs)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeService_Update(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		employeeRepo *mocks.EmployeeRepository
	}

	type args struct {
		employee *entity.Employee
	}

	type expected struct {
		err error
	}

	existingEmployee := &entity.Employee{
		ID:             1,
		BusinessID:     1,
		Specialization: stringPtr("Initial Spec"),
		IsActive:       true,
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: employee successfully updated",
			mock: func(m mocksForExecution) {
				m.employeeRepo.On("Get", ctx, existingEmployee.ID).Return(existingEmployee, nil)
				m.employeeRepo.On("Update", ctx, existingEmployee).Return(nil)
			},
			args: args{
				employee: &entity.Employee{
					ID:             existingEmployee.ID,
					Specialization: stringPtr("New Spec"),
					IsActive:       false,
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "negative: employee not found",
			mock: func(m mocksForExecution) {
				m.employeeRepo.On("Get", ctx, 999).Return(nil, repository.ErrNotFound)
			},
			args: args{
				employee: &entity.Employee{
					ID: 999,
				},
			},
			expected: expected{
				err: fmt.Errorf("failed to get existing employee: %w", repository.ErrNotFound),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			employeeRepoMock := mocks.NewEmployeeRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				employeeRepo: employeeRepoMock,
			})

			// Init service
			employeeService := services.NewEmployeeService(&repository.Repositories{
				Employee: employeeRepoMock,
			})

			// Execute
			err := employeeService.Update(ctx, tc.args.employee)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeService_List(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
		employeeRepo *mocks.EmployeeRepository
	}

	type args struct {
		businessID int
	}

	type expected struct {
		employees []entity.Employee
		err       error
	}

	businessID := 1
	employees := []entity.Employee{
		{
			ID:         1,
			BusinessID: businessID,
			IsActive:   true,
		},
		{
			ID:         2,
			BusinessID: businessID,
			IsActive:   true,
		},
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: employees successfully listed",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.employeeRepo.On("List", ctx, businessID).Return(employees, nil)
			},
			args: args{
				businessID: businessID,
			},
			expected: expected{
				employees: employees,
			},
		},
		{
			name: "negative: business not found",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(nil, repository.ErrNotFound)
			},
			args: args{
				businessID: businessID,
			},
			expected: expected{
				err: fmt.Errorf("invalid business: %w", repository.ErrNotFound),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			businessRepoMock := mocks.NewBusinessRepository(t)
			employeeRepoMock := mocks.NewEmployeeRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
				employeeRepo: employeeRepoMock,
			})

			// Init service
			employeeService := services.NewEmployeeService(&repository.Repositories{
				Business: businessRepoMock,
				Employee: employeeRepoMock,
			})

			// Execute
			got, err := employeeService.List(ctx, tc.args.businessID)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.employees, got)
			}
		})
	}
}
