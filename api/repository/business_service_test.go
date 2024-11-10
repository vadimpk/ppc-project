package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

func TestBusinessServiceRepository_Create(t *testing.T) {

	testCases := []struct {
		name    string
		input   *entity.BusinessService
		wantErr error
	}{
		{
			name: "should create service successfully",
			input: &entity.BusinessService{
				BusinessID: businessID,
				Name:       "Haircut",
				Duration:   30,
				Price:      1000,
				IsActive:   true,
			},
		},
		{
			name: "should create service with description",
			input: &entity.BusinessService{
				BusinessID:  businessID,
				Name:        "Styling",
				Description: stringPtr("Professional hair styling"),
				Duration:    60,
				Price:       2000,
				IsActive:    true,
			},
		},
		{
			name: "should fail with invalid business ID",
			input: &entity.BusinessService{
				BusinessID: 99999,
				Name:       "Invalid Service",
				Duration:   30,
				Price:      1000,
			},
			wantErr: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			err := serviceRepo.Create(context.Background(), tc.input)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tc.input.ID)
			assert.NotZero(t, tc.input.CreatedAt)

			t.Cleanup(func() {
				_ = serviceRepo.Delete(context.Background(), tc.input.ID)
			})

			// Verify created service
			got, err := serviceRepo.Get(context.Background(), tc.input.ID)
			require.NoError(t, err)
			assert.Equal(t, tc.input.Name, got.Name)
			assert.Equal(t, tc.input.Description, got.Description)
			assert.Equal(t, tc.input.Duration, got.Duration)
			assert.Equal(t, tc.input.Price, got.Price)
			assert.Equal(t, tc.input.IsActive, got.IsActive)
		})
	}
}

func TestBusinessServiceRepository_Get(t *testing.T) {

	testCases := []struct {
		name         string
		inputService *entity.BusinessService
		inputID      int
		want         *entity.BusinessService
		wantErr      error
	}{
		{
			name: "should get service by id",
			inputService: &entity.BusinessService{
				BusinessID: businessID,
				Name:       "Test Service",
				Duration:   30,
				Price:      1000,
				IsActive:   true,
			},
			inputID: 0, // will be set after creation
			want: &entity.BusinessService{
				BusinessID: businessID,
				Name:       "Test Service",
				Duration:   30,
				Price:      1000,
				IsActive:   true,
			},
		},
		{
			name: "should get service with description",
			inputService: &entity.BusinessService{
				BusinessID:  businessID,
				Name:        "Test Service with Description",
				Description: stringPtr("Test description"),
				Duration:    45,
				Price:       1500,
				IsActive:    true,
			},
			inputID: 0, // will be set after creation
			want: &entity.BusinessService{
				BusinessID:  businessID,
				Name:        "Test Service with Description",
				Description: stringPtr("Test description"),
				Duration:    45,
				Price:       1500,
				IsActive:    true,
			},
		},
		{
			name:    "should return error when service not found",
			inputID: 99999,
			wantErr: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.inputService != nil {
				err := serviceRepo.Create(context.Background(), tc.inputService)
				require.NoError(t, err)
				tc.inputID = tc.inputService.ID
				tc.want.ID = tc.inputService.ID
				tc.want.CreatedAt = tc.inputService.CreatedAt
			}

			t.Cleanup(func() {
				if tc.inputService != nil {
					_ = serviceRepo.Delete(context.Background(), tc.inputService.ID)
				}
			})

			got, err := serviceRepo.Get(context.Background(), tc.inputID)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want.ID, got.ID)
			assert.Equal(t, tc.want.BusinessID, got.BusinessID)
			assert.Equal(t, tc.want.Name, got.Name)
			assert.Equal(t, tc.want.Description, got.Description)
			assert.Equal(t, tc.want.Duration, got.Duration)
			assert.Equal(t, tc.want.Price, got.Price)
			assert.Equal(t, tc.want.IsActive, got.IsActive)
			assert.Equal(t, tc.want.CreatedAt, got.CreatedAt)
		})
	}
}

func TestBusinessServiceRepository_ListActive(t *testing.T) {

	testServices := []*entity.BusinessService{
		{
			BusinessID: businessID,
			Name:       "Active Service 1",
			Duration:   30,
			Price:      1000,
			IsActive:   true,
		},
		{
			BusinessID: businessID,
			Name:       "Inactive Service",
			Duration:   45,
			Price:      1500,
			IsActive:   false,
		},
		{
			BusinessID: businessID,
			Name:       "Active Service 2",
			Duration:   60,
			Price:      2000,
			IsActive:   true,
		},
	}

	testCases := []struct {
		name          string
		setupServices []*entity.BusinessService
		inputID       int
		wantCount     int
		wantErr       error
	}{
		{
			name:          "should list only active services",
			setupServices: testServices,
			inputID:       businessID,
			wantCount:     2, // only active services
		},
		{
			name:          "should return empty list for non-existent business",
			setupServices: testServices,
			inputID:       99999,
			wantCount:     0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			// Setup services
			var createdServices []int
			for _, s := range tc.setupServices {
				err := serviceRepo.Create(context.Background(), s)
				require.NoError(t, err)
				createdServices = append(createdServices, s.ID)
			}

			t.Cleanup(func() {
				for _, id := range createdServices {
					_ = serviceRepo.Delete(context.Background(), id)
				}
			})

			got, err := serviceRepo.ListActive(context.Background(), tc.inputID)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantCount)

			if tc.wantCount > 0 {
				// Verify all services are active and belong to the correct business
				for _, s := range got {
					assert.True(t, s.IsActive)
					assert.Equal(t, tc.inputID, s.BusinessID)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
