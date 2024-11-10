package repository_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

var (
	db           *repository.DB
	businessRepo repository.BusinessRepository
	serviceRepo  repository.BusinessServiceRepository
	userRepo     repository.UserRepository
	employeeRepo repository.EmployeeRepository
	businessID   int
)

func TestMain(m *testing.M) {
	// Setup
	var err error
	db, err = repository.NewDB(repository.Options{
		Host:           "localhost",
		Port:           5432,
		User:           "postgres",
		Pass:           "postgres",
		DBName:         "ppc_test",
		MinConnections: 1,
		MaxConnections: 2,
		Timezone:       "UTC",
	})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create a test business for all service tests
	business := &entity.Business{Name: "Test Business"}
	businessRepo = repository.NewBusinessRepository(db)
	err = businessRepo.Create(context.Background(), business)
	if err != nil {
		log.Fatalf("Failed to create test business: %v", err)
	}
	businessID = business.ID

	serviceRepo = repository.NewBusinessServiceRepository(db)
	userRepo = repository.NewUserRepository(db)
	employeeRepo = repository.NewEmployeeRepository(db)

	// Run tests
	code := m.Run()

	// Cleanup
	db.PGX.Exec(context.Background(), "TRUNCATE businesses CASCADE")
	db.Close()

	os.Exit(code)
}

func TestBusinessRepository_Get(t *testing.T) {

	testCases := []struct {
		name          string
		inputBusiness *entity.Business
		inputID       int
		want          *entity.Business
		wantErr       error
	}{
		{
			name: "should return business by id",
			inputBusiness: &entity.Business{
				Name: "Test Business",
			},
			inputID: 0, // will be set after creation
			want: &entity.Business{
				Name: "Test Business",
			},
		},
		{
			name: "should return business with full data",
			inputBusiness: &entity.Business{
				Name:    "Test Business Full",
				LogoURL: stringPtr("https://example.com/logo.png"),
				ColorScheme: map[string]interface{}{
					"primary": "#FF0000",
				},
			},
			inputID: 0, // will be set after creation
			want: &entity.Business{
				Name:    "Test Business Full",
				LogoURL: stringPtr("https://example.com/logo.png"),
				ColorScheme: map[string]interface{}{
					"primary": "#FF0000",
				},
			},
		},
		{
			name:    "should return error when business not found",
			inputID: 99999,
			want:    nil,
			wantErr: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.inputBusiness != nil {
				err := businessRepo.Create(context.Background(), tc.inputBusiness)
				require.NoError(t, err)
				tc.inputID = tc.inputBusiness.ID
				tc.want.ID = tc.inputBusiness.ID
				tc.want.CreatedAt = tc.inputBusiness.CreatedAt
			}

			t.Cleanup(func() {
				if tc.inputBusiness != nil {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM businesses WHERE id = $1", tc.inputBusiness.ID)
					require.NoError(t, err)
				}
			})

			got, err := businessRepo.Get(context.Background(), tc.inputID)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want.ID, got.ID)
				assert.Equal(t, tc.want.Name, got.Name)
				assert.Equal(t, tc.want.LogoURL, got.LogoURL)
				assert.Equal(t, tc.want.ColorScheme, got.ColorScheme)
				assert.Equal(t, tc.want.CreatedAt, got.CreatedAt)
			}
		})
	}
}

func TestBusinessRepository_Create(t *testing.T) {

	testCases := []struct {
		name    string
		input   *entity.Business
		want    *entity.Business
		wantErr bool
	}{
		{
			name: "should create business successfully",
			input: &entity.Business{
				Name: "New Business",
			},
			want: &entity.Business{
				Name: "New Business",
			},
		},
		{
			name: "should create business with full data",
			input: &entity.Business{
				Name:    "New Business Full",
				LogoURL: stringPtr("https://example.com/logo.png"),
				ColorScheme: map[string]interface{}{
					"primary":   "#FF0000",
					"secondary": "#00FF00",
				},
			},
			want: &entity.Business{
				Name:    "New Business Full",
				LogoURL: stringPtr("https://example.com/logo.png"),
				ColorScheme: map[string]interface{}{
					"primary":   "#FF0000",
					"secondary": "#00FF00",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			t.Cleanup(func() {
				if tc.input.ID != 0 {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM businesses WHERE id = $1", tc.input.ID)
					require.NoError(t, err)
				}
			})

			err := businessRepo.Create(context.Background(), tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tc.input.ID)
			assert.NotZero(t, tc.input.CreatedAt)

			// Verify created business
			got, err := businessRepo.Get(context.Background(), tc.input.ID)
			require.NoError(t, err)
			assert.Equal(t, tc.input.Name, got.Name)
			assert.Equal(t, tc.input.LogoURL, got.LogoURL)
			assert.Equal(t, tc.input.ColorScheme, got.ColorScheme)
		})
	}
}

func TestBusinessRepository_Update(t *testing.T) {

	testCases := []struct {
		name          string
		inputBusiness *entity.Business
		updateWith    *entity.Business
		want          *entity.Business
		wantErr       error
	}{
		{
			name: "should update business name",
			inputBusiness: &entity.Business{
				Name: "Original Name",
			},
			updateWith: &entity.Business{
				Name: "Updated Name",
			},
			want: &entity.Business{
				Name: "Updated Name",
			},
		},
		{
			name: "should fail when business not found",
			updateWith: &entity.Business{
				ID:   99999,
				Name: "Updated Name",
			},
			wantErr: repository.ErrNotFound,
		},
		{
			name: "should fail with empty name",
			inputBusiness: &entity.Business{
				Name: "Original Name",
			},
			updateWith: &entity.Business{
				Name: "",
			},
			want: &entity.Business{
				Name: "",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.inputBusiness != nil {
				err := businessRepo.Create(context.Background(), tc.inputBusiness)
				require.NoError(t, err)
				tc.updateWith.ID = tc.inputBusiness.ID
				if tc.want != nil {
					tc.want.ID = tc.inputBusiness.ID
				}
			}

			t.Cleanup(func() {
				if tc.inputBusiness != nil {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM businesses WHERE id = $1", tc.inputBusiness.ID)
					require.NoError(t, err)
				}
			})

			err := businessRepo.Update(context.Background(), tc.updateWith)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)

			// Verify updated business
			got, err := businessRepo.Get(context.Background(), tc.updateWith.ID)
			require.NoError(t, err)
			assert.Equal(t, tc.want.Name, got.Name)
		})
	}
}

func TestBusinessRepository_UpdateAppearance(t *testing.T) {

	testCases := []struct {
		name          string
		inputBusiness *entity.Business
		inputLogo     string
		inputColors   map[string]interface{}
		want          *entity.Business
		wantErr       error
	}{
		{
			name: "should update appearance successfully",
			inputBusiness: &entity.Business{
				Name: "Test Business",
			},
			inputLogo: "https://example.com/newlogo.png",
			inputColors: map[string]interface{}{
				"primary": "#0000FF",
			},
			want: &entity.Business{
				Name:    "Test Business",
				LogoURL: stringPtr("https://example.com/newlogo.png"),
				ColorScheme: map[string]interface{}{
					"primary": "#0000FF",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.inputBusiness != nil {
				err := businessRepo.Create(context.Background(), tc.inputBusiness)
				require.NoError(t, err)
				tc.want.ID = tc.inputBusiness.ID
			}

			t.Cleanup(func() {
				if tc.inputBusiness != nil {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM businesses WHERE id = $1", tc.inputBusiness.ID)
					require.NoError(t, err)
				}
			})

			err := businessRepo.UpdateAppearance(context.Background(), tc.inputBusiness.ID, tc.inputLogo, tc.inputColors)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)

			// Verify updated appearance
			got, err := businessRepo.Get(context.Background(), tc.inputBusiness.ID)
			require.NoError(t, err)
			assert.Equal(t, tc.want.LogoURL, got.LogoURL)
			assert.Equal(t, tc.want.ColorScheme, got.ColorScheme)
		})
	}
}
