package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository"
)

func TestUserRepository_Create(t *testing.T) {

	testCases := []struct {
		name         string
		input        *entity.User
		precondition func(t *testing.T)
		validate     func(t *testing.T, got *entity.User)
		wantErr      error
	}{
		{
			name: "should create client with email",
			input: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("test@example.com"),
				FullName:     "Test User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			validate: func(t *testing.T, got *entity.User) {
				assert.Equal(t, "test@example.com", *got.Email)
				assert.Nil(t, got.Phone)
				assert.Equal(t, entity.RoleClient, got.Role)
			},
		},
		{
			name: "should create client with phone",
			input: &entity.User{
				BusinessID:   businessID,
				Phone:        stringPtr("+1234567890"),
				FullName:     "Test User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			validate: func(t *testing.T, got *entity.User) {
				assert.Equal(t, "+1234567890", *got.Phone)
				assert.Nil(t, got.Email)
			},
		},
		{
			name: "should create client with both email and phone",
			input: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("test@example.com"),
				Phone:        stringPtr("+1234567890"),
				FullName:     "Test User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			validate: func(t *testing.T, got *entity.User) {
				assert.Equal(t, "test@example.com", *got.Email)
				assert.Equal(t, "+1234567890", *got.Phone)
			},
		},
		{
			name: "should fail with duplicate email",
			input: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("duplicate@example.com"),
				FullName:     "Test User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			precondition: func(t *testing.T) {
				// First create a user with the same email
				user := &entity.User{
					BusinessID:   businessID,
					Email:        stringPtr("duplicate@example.com"),
					FullName:     "First User",
					PasswordHash: "hash",
					Role:         entity.RoleClient,
				}
				err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
			},
			wantErr: repository.ErrAlreadyExists,
		},
		{
			name: "should fail with duplicate phone",
			input: &entity.User{
				BusinessID:   businessID,
				Phone:        stringPtr("+9999999999"),
				FullName:     "Test User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			precondition: func(t *testing.T) {
				// First create a user with the same phone
				user := &entity.User{
					BusinessID:   businessID,
					Phone:        stringPtr("+9999999999"),
					FullName:     "First User",
					PasswordHash: "hash",
					Role:         entity.RoleClient,
				}
				err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
			},
			wantErr: repository.ErrAlreadyExists,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.precondition != nil {
				tc.precondition(t)
			}

			err := userRepo.Create(context.Background(), tc.input)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tc.input.ID)
			assert.NotZero(t, tc.input.CreatedAt)

			t.Cleanup(func() {
				_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", tc.input.ID)
				require.NoError(t, err)
			})

			// Verify created user
			got, err := userRepo.Get(context.Background(), tc.input.ID)
			require.NoError(t, err)
			assert.Equal(t, tc.input.BusinessID, got.BusinessID)
			assert.Equal(t, tc.input.FullName, got.FullName)
			assert.Equal(t, tc.input.Role, got.Role)
			tc.validate(t, got)
		})
	}
}

func TestUserRepository_CreateBusinessAdmin(t *testing.T) {

	testCases := []struct {
		name         string
		businessName string
		input        *entity.User
		validate     func(t *testing.T, got *entity.User)
		wantErr      error
	}{
		{
			name:         "should create business admin with email",
			businessName: "Test Business",
			input: &entity.User{
				Email:        stringPtr("admin@example.com"),
				FullName:     "Admin User",
				PasswordHash: "hash",
				Role:         entity.RoleAdmin,
			},
			validate: func(t *testing.T, got *entity.User) {
				assert.Equal(t, "admin@example.com", *got.Email)
				assert.Equal(t, entity.RoleAdmin, got.Role)
				assert.NotZero(t, got.BusinessID)
			},
		},
		{
			name:         "should create business admin with phone",
			businessName: "Test Business 2",
			input: &entity.User{
				Phone:        stringPtr("+1234567890"),
				FullName:     "Admin User",
				PasswordHash: "hash",
				Role:         entity.RoleAdmin,
			},
			validate: func(t *testing.T, got *entity.User) {
				assert.Equal(t, "+1234567890", *got.Phone)
				assert.NotZero(t, got.BusinessID)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			err := userRepo.CreateBusinessAdmin(context.Background(), tc.businessName, tc.input)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tc.input.ID)
			assert.NotZero(t, tc.input.BusinessID)
			assert.NotZero(t, tc.input.CreatedAt)

			t.Cleanup(func() {
				_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", tc.input.ID)
				require.NoError(t, err)
				_, err = db.PGX.Exec(context.Background(), "DELETE FROM businesses WHERE id = $1", tc.input.BusinessID)
				require.NoError(t, err)
			})

			// Verify created admin
			got, err := userRepo.Get(context.Background(), tc.input.ID)
			require.NoError(t, err)
			assert.Equal(t, tc.input.FullName, got.FullName)
			tc.validate(t, got)
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {

	testCases := []struct {
		name       string
		inputUser  *entity.User
		inputEmail string
		wantErr    error
	}{
		{
			name: "should get user by email",
			inputUser: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("get@example.com"),
				FullName:     "Get User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			inputEmail: "get@example.com",
		},
		{
			name: "should return not found for non-existent email",
			inputUser: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("exists@example.com"),
				FullName:     "Existing User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			inputEmail: "notfound@example.com",
			wantErr:    repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.inputUser != nil {
				err := userRepo.Create(context.Background(), tc.inputUser)
				require.NoError(t, err)
			}

			t.Cleanup(func() {
				if tc.inputUser != nil {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", tc.inputUser.ID)
					require.NoError(t, err)
				}
			})

			got, err := userRepo.GetByEmail(context.Background(), businessID, tc.inputEmail)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.inputEmail, *got.Email)
			assert.Equal(t, tc.inputUser.BusinessID, got.BusinessID)
		})
	}
}

func TestUserRepository_GetByPhone(t *testing.T) {

	testCases := []struct {
		name       string
		inputUser  *entity.User
		inputPhone string
		wantErr    error
	}{
		{
			name: "should get user by phone",
			inputUser: &entity.User{
				BusinessID:   businessID,
				Phone:        stringPtr("+1111111111"),
				FullName:     "Phone User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			inputPhone: "+1111111111",
		},
		{
			name: "should return not found for non-existent phone",
			inputUser: &entity.User{
				BusinessID:   businessID,
				Phone:        stringPtr("+2222222222"),
				FullName:     "Existing User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			inputPhone: "+3333333333",
			wantErr:    repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.inputUser != nil {
				err := userRepo.Create(context.Background(), tc.inputUser)
				require.NoError(t, err)
			}

			t.Cleanup(func() {
				if tc.inputUser != nil {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", tc.inputUser.ID)
					require.NoError(t, err)
				}
			})

			got, err := userRepo.GetByPhone(context.Background(), businessID, tc.inputPhone)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.inputPhone, *got.Phone)
			assert.Equal(t, tc.inputUser.BusinessID, got.BusinessID)
		})
	}
}

func TestUserRepository_Update(t *testing.T) {

	testCases := []struct {
		name         string
		setupUser    *entity.User
		updateWith   *entity.User
		validate     func(t *testing.T, got *entity.User)
		precondition func(t *testing.T)
		wantErr      error
	}{
		{
			name: "should update user full name",
			setupUser: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("update@example.com"),
				FullName:     "Original Name",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			updateWith: &entity.User{
				FullName: "Updated Name",
			},
			validate: func(t *testing.T, got *entity.User) {
				assert.Equal(t, "Updated Name", got.FullName)
				assert.Equal(t, "update@example.com", *got.Email)
			},
		},
		{
			name: "should update user email",
			setupUser: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("old@example.com"),
				FullName:     "Test User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			updateWith: &entity.User{
				Email: stringPtr("new@example.com"),
			},
			validate: func(t *testing.T, got *entity.User) {
				assert.Equal(t, "new@example.com", *got.Email)
			},
		},
		{
			name: "should fail with duplicate email",
			setupUser: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("first@example.com"),
				FullName:     "First User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			updateWith: &entity.User{
				Email: stringPtr("taken@example.com"),
			},
			precondition: func(t *testing.T) {
				// Create another user with the email we'll try to update to
				other := &entity.User{
					BusinessID:   businessID,
					Email:        stringPtr("taken@example.com"),
					FullName:     "Other User",
					PasswordHash: "hash",
					Role:         entity.RoleClient,
				}
				err := userRepo.Create(context.Background(), other)
				require.NoError(t, err)
			},
			wantErr: repository.ErrAlreadyExists,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			// Setup initial user
			err := userRepo.Create(context.Background(), tc.setupUser)
			require.NoError(t, err)

			tc.updateWith.ID = tc.setupUser.ID

			if tc.precondition != nil {
				tc.precondition(t)
			}

			t.Cleanup(func() {
				_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE business_id = $1", businessID)
				require.NoError(t, err)
			})

			err = userRepo.Update(context.Background(), tc.updateWith)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)

			// Verify updated user
			got, err := userRepo.Get(context.Background(), tc.setupUser.ID)
			require.NoError(t, err)
			tc.validate(t, got)
		})
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {

	testCases := []struct {
		name          string
		setupUser     *entity.User
		newPassword   string
		checkPassword bool
		wantErr       error
	}{
		{
			name: "should update password",
			setupUser: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("pwd@example.com"),
				FullName:     "Password User",
				PasswordHash: "original_hash",
				Role:         entity.RoleClient,
			},
			newPassword:   "new_hash",
			checkPassword: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			// Setup user if needed
			if tc.setupUser.Email != nil {
				err := userRepo.Create(context.Background(), tc.setupUser)
				require.NoError(t, err)
			}

			t.Cleanup(func() {
				if tc.setupUser.Email != nil {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", tc.setupUser.ID)
					require.NoError(t, err)
				}
			})

			err := userRepo.UpdatePassword(context.Background(), tc.setupUser.ID, tc.newPassword)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)

			if tc.checkPassword {
				// Verify password was updated
				got, err := userRepo.Get(context.Background(), tc.setupUser.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.newPassword, got.PasswordHash)
			}
		})
	}
}

func TestUserRepository_Get(t *testing.T) {

	testCases := []struct {
		name      string
		setupUser *entity.User
		inputID   int
		wantErr   error
	}{
		{
			name: "should get user by id",
			setupUser: &entity.User{
				BusinessID:   businessID,
				Email:        stringPtr("get@example.com"),
				Phone:        stringPtr("+1234567890"),
				FullName:     "Get User",
				PasswordHash: "hash",
				Role:         entity.RoleClient,
			},
			inputID: 0, // will be set after creation
		},
		{
			name:    "should return error when user not found",
			inputID: 99999,
			wantErr: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			if tc.setupUser != nil {
				err := userRepo.Create(context.Background(), tc.setupUser)
				require.NoError(t, err)
				tc.inputID = tc.setupUser.ID
			}

			t.Cleanup(func() {
				if tc.setupUser != nil {
					_, err := db.PGX.Exec(context.Background(), "DELETE FROM users WHERE id = $1", tc.setupUser.ID)
					require.NoError(t, err)
				}
			})

			got, err := userRepo.Get(context.Background(), tc.inputID)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.setupUser.ID, got.ID)
			assert.Equal(t, tc.setupUser.BusinessID, got.BusinessID)
			assert.Equal(t, tc.setupUser.Email, got.Email)
			assert.Equal(t, tc.setupUser.Phone, got.Phone)
			assert.Equal(t, tc.setupUser.FullName, got.FullName)
			assert.Equal(t, tc.setupUser.Role, got.Role)
		})
	}
}
