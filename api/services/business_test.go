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

func TestBusinessService_Create(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
	}

	type args struct {
		business *entity.Business
	}

	type expected struct {
		err error
	}

	business := &entity.Business{
		Name: "Test Business",
		ColorScheme: map[string]interface{}{
			"primary":    "#FF0000",
			"secondary":  "#00FF00",
			"background": "#0000FF",
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
			name: "positive: business successfully created",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Create", ctx, business).Return(nil)
			},
			args: args{
				business: business,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "negative: empty business name",
			mock: func(m mocksForExecution) {},
			args: args{
				business: &entity.Business{
					Name: "",
				},
			},
			expected: expected{
				err: fmt.Errorf("business name is required"),
			},
		},
		{
			name: "negative: failed to create business",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Create", ctx, business).Return(fmt.Errorf("some error"))
			},
			args: args{
				business: business,
			},
			expected: expected{
				err: fmt.Errorf("failed to create business: %w", fmt.Errorf("some error")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			businessRepoMock := mocks.NewBusinessRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
			})

			// Init service
			businessService := services.NewBusinessService(&repository.Repositories{
				Business: businessRepoMock,
			})

			// Execute
			err := businessService.Create(ctx, tc.args.business)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBusinessService_Get(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
	}

	type args struct {
		id int
	}

	type expected struct {
		business *entity.Business
		err      error
	}

	business := &entity.Business{
		ID:   1,
		Name: "Test Business",
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: business successfully retrieved",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, business.ID).Return(business, nil)
			},
			args: args{
				id: business.ID,
			},
			expected: expected{
				business: business,
				err:      nil,
			},
		},
		{
			name: "negative: business not found",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, 999).Return(nil, repository.ErrNotFound)
			},
			args: args{
				id: 999,
			},
			expected: expected{
				business: nil,
				err:      fmt.Errorf("failed to get business: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: repository error",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, business.ID).Return(nil, fmt.Errorf("some error"))
			},
			args: args{
				id: business.ID,
			},
			expected: expected{
				business: nil,
				err:      fmt.Errorf("failed to get business: %w", fmt.Errorf("some error")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			businessRepoMock := mocks.NewBusinessRepository(t)
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
			})

			businessService := services.NewBusinessService(&repository.Repositories{
				Business: businessRepoMock,
			})

			got, err := businessService.Get(ctx, tc.args.id)

			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.business, got)
			}
		})
	}
}

func TestBusinessService_Update(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
	}

	type args struct {
		business *entity.Business
	}

	type expected struct {
		err error
	}

	existingBusiness := &entity.Business{
		ID:   1,
		Name: "Original Name",
	}

	updatedBusiness := &entity.Business{
		ID:   1,
		Name: "Updated Name",
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: business successfully updated",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, existingBusiness.ID).Return(existingBusiness, nil)
				m.businessRepo.On("Update", ctx, existingBusiness).Return(nil)
			},
			args: args{
				business: updatedBusiness,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "negative: business not found",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, updatedBusiness.ID).Return(nil, repository.ErrNotFound)
			},
			args: args{
				business: updatedBusiness,
			},
			expected: expected{
				err: fmt.Errorf("failed to get existing business: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: empty business name",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, existingBusiness.ID).Return(existingBusiness, nil)
			},
			args: args{
				business: &entity.Business{
					ID:   1,
					Name: "",
				},
			},
			expected: expected{
				err: fmt.Errorf("business name is required"),
			},
		},
		{
			name: "negative: failed to update business",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, existingBusiness.ID).Return(existingBusiness, nil)
				m.businessRepo.On("Update", ctx, existingBusiness).Return(fmt.Errorf("some error"))
			},
			args: args{
				business: updatedBusiness,
			},
			expected: expected{
				err: fmt.Errorf("failed to update business: %w", fmt.Errorf("some error")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			businessRepoMock := mocks.NewBusinessRepository(t)
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
			})

			businessService := services.NewBusinessService(&repository.Repositories{
				Business: businessRepoMock,
			})

			err := businessService.Update(ctx, tc.args.business)

			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBusinessService_UpdateAppearance(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
	}

	type args struct {
		id          int
		logoURL     string
		colorScheme map[string]interface{}
	}

	type expected struct {
		err error
	}

	validColorScheme := map[string]interface{}{
		"primary":    "#FF0000",
		"secondary":  "#00FF00",
		"background": "#0000FF",
	}

	ctx := context.Background()
	businessID := 1

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: appearance successfully updated",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.businessRepo.On("UpdateAppearance", ctx, businessID, "logo.png", validColorScheme).Return(nil)
			},
			args: args{
				id:          businessID,
				logoURL:     "logo.png",
				colorScheme: validColorScheme,
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
				id:          businessID,
				logoURL:     "logo.png",
				colorScheme: validColorScheme,
			},
			expected: expected{
				err: fmt.Errorf("failed to get existing business: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: invalid color scheme",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
			},
			args: args{
				id:      businessID,
				logoURL: "logo.png",
				colorScheme: map[string]interface{}{
					"primary": "#FF0000",
					// missing required colors
				},
			},
			expected: expected{
				err: fmt.Errorf("invalid color scheme: %w", fmt.Errorf("missing required color: secondary")),
			},
		},
		{
			name: "negative: failed to update appearance",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.businessRepo.On("UpdateAppearance", ctx, businessID, "logo.png", validColorScheme).Return(fmt.Errorf("some error"))
			},
			args: args{
				id:          businessID,
				logoURL:     "logo.png",
				colorScheme: validColorScheme,
			},
			expected: expected{
				err: fmt.Errorf("failed to update business appearance: %w", fmt.Errorf("some error")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			businessRepoMock := mocks.NewBusinessRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
			})

			// Init service
			businessService := services.NewBusinessService(&repository.Repositories{
				Business: businessRepoMock,
			})

			// Execute
			err := businessService.UpdateAppearance(ctx, tc.args.id, tc.args.logoURL, tc.args.colorScheme)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
