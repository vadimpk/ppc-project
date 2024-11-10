package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

// Helper function to create a test user
func createTestUser(t *testing.T, role string) *entity.User {
	user := &entity.User{
		BusinessID:   businessID,
		Email:        stringPtr("employee@example.com"),
		FullName:     "Test Employee",
		PasswordHash: "hash",
		Role:         role,
	}
	err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)
	return user
}

// Helper function to create a test service
func createTestService(t *testing.T) *entity.BusinessService {
	service := &entity.BusinessService{
		BusinessID: businessID,
		Name:       "Test Service",
		Duration:   30,
		Price:      1000,
		IsActive:   true,
	}
	err := serviceRepo.Create(context.Background(), service)
	require.NoError(t, err)
	return service
}

func TestEmployeeRepository_Create(t *testing.T) {

	testCases := []struct {
		name      string
		input     *entity.Employee
		setupUser bool
		validate  func(t *testing.T, got *entity.Employee)
		wantErr   error
	}{
		{
			name: "should create employee successfully",
			input: &entity.Employee{
				BusinessID:     businessID,
				IsActive:       true,
				Specialization: stringPtr("Hair Stylist"),
			},
			setupUser: true,
			validate: func(t *testing.T, got *entity.Employee) {
				assert.NotZero(t, got.ID)
				assert.Equal(t, businessID, got.BusinessID)
				assert.Equal(t, "Hair Stylist", *got.Specialization)
				assert.True(t, got.IsActive)
				assert.NotNil(t, got.User)
				assert.Equal(t, "Test Employee", got.User.FullName)
			},
		},
		{
			name: "should create employee without specialization",
			input: &entity.Employee{
				BusinessID: businessID,
				IsActive:   true,
			},
			setupUser: true,
			validate: func(t *testing.T, got *entity.Employee) {
				assert.Nil(t, got.Specialization)
				assert.True(t, got.IsActive)
			},
		},
		{
			name: "should fail with non-existent user",
			input: &entity.Employee{
				BusinessID: businessID,
				UserID:     99999,
				IsActive:   true,
			},
			setupUser: false,
			wantErr:   repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.setupUser {
				user := createTestUser(t, entity.RoleEmployee)
				tc.input.UserID = user.ID
			}

			err := employeeRepo.Create(context.Background(), tc.input)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tc.input.ID)
			assert.NotZero(t, tc.input.CreatedAt)

			t.Cleanup(func() {
				_, err := db.PGX.Exec(context.Background(), "DELETE FROM employees WHERE id = $1", tc.input.ID)
				require.NoError(t, err)
				if tc.setupUser {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", tc.input.UserID)
					require.NoError(t, err)
				}
			})

			// Verify created employee
			got, err := employeeRepo.Get(context.Background(), tc.input.ID)
			require.NoError(t, err)
			tc.validate(t, got)
		})
	}
}

func TestEmployeeRepository_Update(t *testing.T) {

	testCases := []struct {
		name          string
		setupEmployee *entity.Employee
		updateWith    *entity.Employee
		validate      func(t *testing.T, got *entity.Employee)
		wantErr       error
	}{
		{
			name: "should update employee specialization",
			setupEmployee: &entity.Employee{
				BusinessID:     businessID,
				Specialization: stringPtr("Initial Spec"),
				IsActive:       true,
			},
			updateWith: &entity.Employee{
				Specialization: stringPtr("Updated Spec"),
				IsActive:       true,
			},
			validate: func(t *testing.T, got *entity.Employee) {
				assert.Equal(t, "Updated Spec", *got.Specialization)
				assert.True(t, got.IsActive)
			},
		},
		{
			name: "should update active status",
			setupEmployee: &entity.Employee{
				BusinessID:     businessID,
				Specialization: stringPtr("Spec"),
				IsActive:       true,
			},
			updateWith: &entity.Employee{
				Specialization: stringPtr("Spec"),
				IsActive:       false,
			},
			validate: func(t *testing.T, got *entity.Employee) {
				assert.False(t, got.IsActive)
			},
		},
		{
			name: "should fail with non-existent employee",
			updateWith: &entity.Employee{
				ID:       99999,
				IsActive: true,
			},
			wantErr: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.setupEmployee != nil {
				user := createTestUser(t, entity.RoleEmployee)
				tc.setupEmployee.UserID = user.ID
				err := employeeRepo.Create(context.Background(), tc.setupEmployee)
				require.NoError(t, err)
				tc.updateWith.ID = tc.setupEmployee.ID
			}

			err := employeeRepo.Update(context.Background(), tc.updateWith)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)

			t.Cleanup(func() {
				if tc.setupEmployee != nil {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM employees WHERE id = $1", tc.setupEmployee.ID)
					require.NoError(t, err)
					_, err = db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", tc.setupEmployee.UserID)
					require.NoError(t, err)
				}
			})

			// Verify updated employee
			got, err := employeeRepo.Get(context.Background(), tc.updateWith.ID)
			require.NoError(t, err)
			tc.validate(t, got)
		})
	}
}

func TestEmployeeRepository_List(t *testing.T) {

	testCases := []struct {
		name            string
		setupEmployees  int
		inputBusinessID int
		wantCount       int
		wantErr         error
	}{
		{
			name:            "should list all employees",
			setupEmployees:  3,
			inputBusinessID: businessID,
			wantCount:       3,
		},
		{
			name:            "should return empty list for non-existent business",
			setupEmployees:  2,
			inputBusinessID: 99999,
			wantCount:       0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			var createdEmployees []int
			var createdUsers []int

			// Setup employees
			for i := 0; i < tc.setupEmployees; i++ {
				user := &entity.User{
					BusinessID:   businessID,
					Email:        stringPtr(fmt.Sprintf("employee%d@example.com", i)),
					FullName:     fmt.Sprintf("Test Employee %d", i),
					PasswordHash: "hash",
					Role:         entity.RoleEmployee,
				}
				err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)

				employee := &entity.Employee{
					BusinessID:     businessID,
					UserID:         user.ID,
					IsActive:       true,
					Specialization: stringPtr(fmt.Sprintf("Spec %d", i)),
				}
				err = employeeRepo.Create(context.Background(), employee)
				require.NoError(t, err)

				createdEmployees = append(createdEmployees, employee.ID)
				createdUsers = append(createdUsers, user.ID)
			}

			t.Cleanup(func() {
				for _, id := range createdEmployees {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM employees WHERE id = $1", id)
					require.NoError(t, err)
				}
				for _, id := range createdUsers {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
					require.NoError(t, err)
				}
			})

			got, err := employeeRepo.List(context.Background(), tc.inputBusinessID)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantCount)

			if tc.wantCount > 0 {
				for _, emp := range got {
					assert.Equal(t, tc.inputBusinessID, emp.BusinessID)
					assert.NotNil(t, emp.User)
				}
			}
		})
	}
}

func TestEmployeeRepository_Services(t *testing.T) {

	testCases := []struct {
		name           string
		setupServices  int
		assignServices []int
		validate       func(t *testing.T, employee *entity.Employee, services []entity.BusinessService)
		wantErr        error
	}{
		{
			name:           "should assign and get services",
			setupServices:  3,
			assignServices: []int{0, 1}, // will be replaced with actual IDs
			validate: func(t *testing.T, employee *entity.Employee, services []entity.BusinessService) {
				assert.Len(t, services, 2)
				assert.Equal(t, "Test Service", services[0].Name)
			},
		},
		{
			name:           "should assign and remove services",
			setupServices:  3,
			assignServices: []int{0, 1, 2}, // will be replaced with actual IDs
			validate: func(t *testing.T, employee *entity.Employee, services []entity.BusinessService) {
				// First verify all services are assigned
				assert.Len(t, services, 3)

				// Remove one service
				err := employeeRepo.RemoveServices(context.Background(), employee.ID, []int{services[0].ID})
				require.NoError(t, err)

				// Verify services after removal
				updatedServices, err := employeeRepo.GetServices(context.Background(), employee.ID)
				require.NoError(t, err)
				assert.Len(t, updatedServices, 2)
			},
		},
		{
			name:           "should handle empty service list",
			setupServices:  0,
			assignServices: []int{},
			validate: func(t *testing.T, employee *entity.Employee, services []entity.BusinessService) {
				assert.Len(t, services, 0)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			// Create employee
			user := createTestUser(t, entity.RoleEmployee)
			employee := &entity.Employee{
				BusinessID: businessID,
				UserID:     user.ID,
				IsActive:   true,
			}
			err := employeeRepo.Create(context.Background(), employee)
			require.NoError(t, err)

			// Create services
			var createdServices []int
			for i := 0; i < tc.setupServices; i++ {
				service := createTestService(t)
				createdServices = append(createdServices, service.ID)
			}

			// Replace placeholder IDs with actual service IDs
			serviceIDs := make([]int, len(tc.assignServices))
			for i, idx := range tc.assignServices {
				serviceIDs[i] = createdServices[idx]
			}

			t.Cleanup(func() {
				_, err := db.PGX.Exec(context.Background(), "DELETE FROM employee_services WHERE employee_id = $1", employee.ID)
				require.NoError(t, err)
				_, err = db.PGX.Exec(context.Background(), "DELETE FROM employees WHERE id = $1", employee.ID)
				require.NoError(t, err)
				_, err = db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", user.ID)
				require.NoError(t, err)
				for _, id := range createdServices {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM services WHERE id = $1", id)
					require.NoError(t, err)
				}
			})

			// Assign services
			if len(serviceIDs) > 0 {
				err = employeeRepo.AssignServices(context.Background(), employee.ID, serviceIDs)
				if tc.wantErr != nil {
					assert.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
			}

			// Get services
			services, err := employeeRepo.GetServices(context.Background(), employee.ID)
			require.NoError(t, err)

			// Validate
			tc.validate(t, employee, services)
		})
	}
}
