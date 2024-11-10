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

func TestBusinessServiceService_Create(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
		serviceRepo  *mocks.BusinessServiceRepository
	}

	type args struct {
		service *entity.BusinessService
	}

	type expected struct {
		err error
	}

	businessID := 1
	service := &entity.BusinessService{
		BusinessID: businessID,
		Name:       "Haircut",
		Duration:   30,
		Price:      1000,
		IsActive:   true,
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: service successfully created",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.serviceRepo.On("Create", ctx, service).Return(nil)
			},
			args: args{
				service: service,
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
				service: service,
			},
			expected: expected{
				err: fmt.Errorf("invalid business: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: empty service name",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
			},
			args: args{
				service: &entity.BusinessService{
					BusinessID: businessID,
					Name:       "",
					Duration:   30,
					Price:      1000,
				},
			},
			expected: expected{
				err: fmt.Errorf("service name is required"),
			},
		},
		{
			name: "negative: invalid duration",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
			},
			args: args{
				service: &entity.BusinessService{
					BusinessID: businessID,
					Name:       "Test Service",
					Duration:   0,
					Price:      1000,
				},
			},
			expected: expected{
				err: fmt.Errorf("service duration must be positive"),
			},
		},
		{
			name: "negative: invalid price",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
			},
			args: args{
				service: &entity.BusinessService{
					BusinessID: businessID,
					Name:       "Test Service",
					Duration:   30,
					Price:      -100,
				},
			},
			expected: expected{
				err: fmt.Errorf("service price cannot be negative"),
			},
		},
		{
			name: "negative: failed to create service",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.serviceRepo.On("Create", ctx, service).Return(fmt.Errorf("some error"))
			},
			args: args{
				service: service,
			},
			expected: expected{
				err: fmt.Errorf("failed to create service: %w", fmt.Errorf("some error")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			businessRepoMock := mocks.NewBusinessRepository(t)
			serviceRepoMock := mocks.NewBusinessServiceRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
				serviceRepo:  serviceRepoMock,
			})

			// Init service
			businessService := services.NewBusinessServiceService(&repository.Repositories{
				Business: businessRepoMock,
				Service:  serviceRepoMock,
			})

			// Execute
			err := businessService.Create(ctx, tc.args.service)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBusinessServiceService_Update(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		serviceRepo *mocks.BusinessServiceRepository
	}

	type args struct {
		service *entity.BusinessService
	}

	type expected struct {
		err error
	}

	existingService := &entity.BusinessService{
		ID:         1,
		BusinessID: 1,
		Name:       "Original Service",
		Duration:   30,
		Price:      1000,
		IsActive:   true,
	}

	updatedService := &entity.BusinessService{
		ID:         1,
		BusinessID: 1,
		Name:       "Updated Service",
		Duration:   45,
		Price:      1500,
		IsActive:   true,
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: service successfully updated",
			mock: func(m mocksForExecution) {
				m.serviceRepo.On("Get", ctx, existingService.ID).Return(existingService, nil)
				m.serviceRepo.On("Update", ctx, updatedService).Return(nil)
			},
			args: args{
				service: updatedService,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "negative: service not found",
			mock: func(m mocksForExecution) {
				m.serviceRepo.On("Get", ctx, updatedService.ID).Return(nil, repository.ErrNotFound)
			},
			args: args{
				service: updatedService,
			},
			expected: expected{
				err: fmt.Errorf("invalid service: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: service does not belong to business",
			mock: func(m mocksForExecution) {
				m.serviceRepo.On("Get", ctx, existingService.ID).Return(existingService, nil)
			},
			args: args{
				service: &entity.BusinessService{
					ID:         existingService.ID,
					BusinessID: 999, // different business ID
					Name:       "Updated Service",
					Duration:   45,
					Price:      1500,
				},
			},
			expected: expected{
				err: fmt.Errorf("service does not belong to the business"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			serviceRepoMock := mocks.NewBusinessServiceRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				serviceRepo: serviceRepoMock,
			})

			// Init service
			businessService := services.NewBusinessServiceService(&repository.Repositories{
				Service: serviceRepoMock,
			})

			// Execute
			err := businessService.Update(ctx, tc.args.service)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBusinessServiceService_List(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
		serviceRepo  *mocks.BusinessServiceRepository
	}

	type args struct {
		businessID int
	}

	type expected struct {
		services []entity.BusinessService
		err      error
	}

	businessID := 1
	businessServices := []entity.BusinessService{
		{
			ID:         1,
			BusinessID: businessID,
			Name:       "Service 1",
			Duration:   30,
			Price:      1000,
			IsActive:   true,
		},
		{
			ID:         2,
			BusinessID: businessID,
			Name:       "Service 2",
			Duration:   45,
			Price:      1500,
			IsActive:   false,
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
			name: "positive: businessServices successfully listed",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.serviceRepo.On("List", ctx, businessID).Return(businessServices, nil)
			},
			args: args{
				businessID: businessID,
			},
			expected: expected{
				services: businessServices,
				err:      nil,
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
				services: nil,
				err:      fmt.Errorf("invalid business: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: failed to list businessServices",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.serviceRepo.On("List", ctx, businessID).Return(nil, fmt.Errorf("some error"))
			},
			args: args{
				businessID: businessID,
			},
			expected: expected{
				services: nil,
				err:      fmt.Errorf("failed to list services: %w", fmt.Errorf("some error")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			businessRepoMock := mocks.NewBusinessRepository(t)
			serviceRepoMock := mocks.NewBusinessServiceRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
				serviceRepo:  serviceRepoMock,
			})

			// Init service
			businessService := services.NewBusinessServiceService(&repository.Repositories{
				Business: businessRepoMock,
				Service:  serviceRepoMock,
			})

			// Execute
			got, err := businessService.List(ctx, tc.args.businessID)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.services, got)
			}
		})
	}
}
