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

func TestUserService_Create(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		businessRepo *mocks.BusinessRepository
		userRepo     *mocks.UserRepository
	}

	type args struct {
		user *entity.User
	}

	type expected struct {
		user *entity.User
		err  error
	}

	businessID := 1
	email := "test@example.com"
	phone := "+1234567890"

	userWithEmail := &entity.User{
		BusinessID:   businessID,
		Email:        &email,
		FullName:     "Test User",
		PasswordHash: "hash",
		Role:         entity.RoleClient,
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: user successfully created with email",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.userRepo.On("GetByEmail", ctx, email).Return(nil, repository.ErrNotFound)
				m.userRepo.On("Create", ctx, userWithEmail).Return(nil)
			},
			args: args{
				user: userWithEmail,
			},
			expected: expected{
				user: userWithEmail,
			},
		},
		{
			name: "positive: admin user created without business check",
			mock: func(m mocksForExecution) {
				adminUser := &entity.User{
					BusinessID:   businessID,
					Email:        &email,
					FullName:     "Admin User",
					PasswordHash: "hash",
					Role:         entity.RoleAdmin,
				}
				m.userRepo.On("GetByEmail", ctx, email).Return(nil, repository.ErrNotFound)
				m.userRepo.On("Create", ctx, adminUser).Return(nil)
			},
			args: args{
				user: &entity.User{
					BusinessID:   businessID,
					Email:        &email,
					FullName:     "Admin User",
					PasswordHash: "hash",
					Role:         entity.RoleAdmin,
				},
			},
			expected: expected{
				user: &entity.User{
					BusinessID:   businessID,
					Email:        &email,
					FullName:     "Admin User",
					PasswordHash: "hash",
					Role:         entity.RoleAdmin,
				},
			},
		},
		{
			name: "negative: business not found",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(nil, repository.ErrNotFound)
			},
			args: args{
				user: userWithEmail,
			},
			expected: expected{
				err: fmt.Errorf("invalid business: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: email already exists",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.userRepo.On("GetByEmail", ctx, email).Return(&entity.User{}, nil)
			},
			args: args{
				user: userWithEmail,
			},
			expected: expected{
				err: fmt.Errorf("email already exists"),
			},
		},
		{
			name: "negative: phone already exists",
			mock: func(m mocksForExecution) {
				m.businessRepo.On("Get", ctx, businessID).Return(&entity.Business{ID: businessID}, nil)
				m.userRepo.On("GetByPhone", ctx, phone).Return(&entity.User{}, nil)
			},
			args: args{
				user: &entity.User{
					BusinessID:   businessID,
					Phone:        &phone,
					FullName:     "Test User",
					PasswordHash: "hash",
					Role:         entity.RoleClient,
				},
			},
			expected: expected{
				err: fmt.Errorf("phone already exists"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			businessRepoMock := mocks.NewBusinessRepository(t)
			userRepoMock := mocks.NewUserRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				businessRepo: businessRepoMock,
				userRepo:     userRepoMock,
			})

			// Init service
			userService := services.NewUserService(&repository.Repositories{
				Business: businessRepoMock,
				User:     userRepoMock,
			})

			// Execute
			got, err := userService.Create(ctx, tc.args.user)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.user.Email, got.Email)
				assert.Equal(t, tc.expected.user.Phone, got.Phone)
				assert.Equal(t, tc.expected.user.FullName, got.FullName)
				assert.Equal(t, tc.expected.user.Role, got.Role)
			}
		})
	}
}

func TestUserService_CreateBusinessAdmin(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		userRepo *mocks.UserRepository
	}

	type args struct {
		businessName string
		user         *entity.User
	}

	type expected struct {
		user *entity.User
		err  error
	}

	email := "admin@example.com"
	businessName := "Test Business"

	adminUser := &entity.User{
		Email:        &email,
		FullName:     "Admin User",
		PasswordHash: "hash",
		Role:         entity.RoleAdmin,
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: business admin successfully created",
			mock: func(m mocksForExecution) {
				m.userRepo.On("GetByEmail", ctx, email).Return(nil, repository.ErrNotFound)
				m.userRepo.On("CreateBusinessAdmin", ctx, businessName, adminUser).Return(nil)
			},
			args: args{
				businessName: businessName,
				user:         adminUser,
			},
			expected: expected{
				user: adminUser,
			},
		},
		{
			name: "negative: email already exists",
			mock: func(m mocksForExecution) {
				m.userRepo.On("GetByEmail", ctx, email).Return(&entity.User{}, nil)
			},
			args: args{
				businessName: businessName,
				user:         adminUser,
			},
			expected: expected{
				err: fmt.Errorf("email already exists"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			userRepoMock := mocks.NewUserRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				userRepo: userRepoMock,
			})

			// Init service
			userService := services.NewUserService(&repository.Repositories{
				User: userRepoMock,
			})

			// Execute
			got, err := userService.CreateBusinessAdmin(ctx, tc.args.businessName, tc.args.user)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.user.Email, got.Email)
				assert.Equal(t, tc.expected.user.FullName, got.FullName)
				assert.Equal(t, tc.expected.user.Role, got.Role)
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		userRepo *mocks.UserRepository
	}

	type args struct {
		user *entity.User
	}

	type expected struct {
		user *entity.User
		err  error
	}

	existingEmail := "existing@example.com"
	newEmail := "new@example.com"
	businessID := 1

	existingUser := &entity.User{
		ID:         1,
		BusinessID: businessID,
		Email:      &existingEmail,
		FullName:   "Existing User",
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: user successfully updated",
			mock: func(m mocksForExecution) {
				updatedUser := &entity.User{
					ID:         existingUser.ID,
					BusinessID: businessID,
					Email:      &newEmail,
					FullName:   "Updated User",
				}
				m.userRepo.On("Get", ctx, existingUser.ID).Return(existingUser, nil)
				m.userRepo.On("GetByEmail", ctx, newEmail).Return(nil, repository.ErrNotFound)
				m.userRepo.On("Update", ctx, updatedUser).Return(nil)
			},
			args: args{
				user: &entity.User{
					ID:         existingUser.ID,
					BusinessID: businessID,
					Email:      &newEmail,
					FullName:   "Updated User",
				},
			},
			expected: expected{
				user: &entity.User{
					ID:         existingUser.ID,
					BusinessID: businessID,
					Email:      &newEmail,
					FullName:   "Updated User",
				},
			},
		},
		{
			name: "negative: user not found",
			mock: func(m mocksForExecution) {
				m.userRepo.On("Get", ctx, 999).Return(nil, repository.ErrNotFound)
			},
			args: args{
				user: &entity.User{
					ID:         999,
					BusinessID: businessID,
					Email:      &newEmail,
				},
			},
			expected: expected{
				err: fmt.Errorf("failed to get existing user: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: email already exists",
			mock: func(m mocksForExecution) {
				m.userRepo.On("Get", ctx, existingUser.ID).Return(existingUser, nil)
				m.userRepo.On("GetByEmail", ctx, newEmail).Return(&entity.User{}, nil)
			},
			args: args{
				user: &entity.User{
					ID:         existingUser.ID,
					BusinessID: businessID,
					Email:      &newEmail,
				},
			},
			expected: expected{
				err: fmt.Errorf("email already exists"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			userRepoMock := mocks.NewUserRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				userRepo: userRepoMock,
			})

			// Init service
			userService := services.NewUserService(&repository.Repositories{
				User: userRepoMock,
			})

			// Execute
			got, err := userService.Update(ctx, tc.args.user)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.user.Email, got.Email)
				assert.Equal(t, tc.expected.user.FullName, got.FullName)
			}
		})
	}
}

func TestUserService_Authenticate(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		userRepo *mocks.UserRepository
	}

	type args struct {
		businessID int
		email      string
		phone      string
		password   string
	}

	type expected struct {
		user *entity.User
		err  error
	}

	businessID := 1
	email := "test@example.com"
	phone := "+1234567890"
	password := "password"

	user := &entity.User{
		ID:           1,
		BusinessID:   businessID,
		Email:        &email,
		PasswordHash: "hashed_password",
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: authenticate with email",
			mock: func(m mocksForExecution) {
				m.userRepo.On("GetByEmail", ctx, email).Return(user, nil)
			},
			args: args{
				businessID: businessID,
				email:      email,
				password:   password,
			},
			expected: expected{
				user: user,
			},
		},
		{
			name: "positive: authenticate with phone",
			mock: func(m mocksForExecution) {
				userWithPhone := &entity.User{
					ID:           1,
					BusinessID:   businessID,
					Phone:        &phone,
					PasswordHash: "hashed_password",
				}
				m.userRepo.On("GetByPhone", ctx, phone).Return(userWithPhone, nil)
			},
			args: args{
				businessID: businessID,
				phone:      phone,
				password:   password,
			},
			expected: expected{
				user: &entity.User{
					ID:           1,
					BusinessID:   businessID,
					Phone:        &phone,
					PasswordHash: "hashed_password",
				},
			},
		},
		{
			name: "negative: no email or phone provided",
			args: args{
				businessID: businessID,
				password:   password,
			},
			expected: expected{
				err: fmt.Errorf("either email or phone is required"),
			},
		},
		{
			name: "negative: user not found",
			mock: func(m mocksForExecution) {
				m.userRepo.On("GetByEmail", ctx, email).Return(nil, repository.ErrNotFound)
			},
			args: args{
				businessID: businessID,
				email:      email,
				password:   password,
			},
			expected: expected{
				err: fmt.Errorf("invalid credentials"),
			},
		},
		{
			name: "negative: repository error",
			mock: func(m mocksForExecution) {
				m.userRepo.On("GetByEmail", ctx, email).Return(nil, fmt.Errorf("some error"))
			},
			args: args{
				businessID: businessID,
				email:      email,
				password:   password,
			},
			expected: expected{
				err: fmt.Errorf("failed to get user: %w", fmt.Errorf("some error")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			userRepoMock := mocks.NewUserRepository(t)

			// Setup mocks
			if tc.mock != nil {
				tc.mock(mocksForExecution{
					userRepo: userRepoMock,
				})
			}

			// Init service
			userService := services.NewUserService(&repository.Repositories{
				User: userRepoMock,
			})

			// Execute
			got, err := userService.Authenticate(ctx, tc.args.email, tc.args.phone, tc.args.password)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.user.ID, got.ID)
				assert.Equal(t, tc.expected.user.BusinessID, got.BusinessID)
				assert.Equal(t, tc.expected.user.Email, got.Email)
				assert.Equal(t, tc.expected.user.Phone, got.Phone)
				assert.Equal(t, tc.expected.user.PasswordHash, got.PasswordHash)
			}
		})
	}
}

func TestUserService_Get(t *testing.T) {
	t.Parallel()

	type mocksForExecution struct {
		userRepo *mocks.UserRepository
	}

	type args struct {
		id int
	}

	type expected struct {
		user *entity.User
		err  error
	}

	user := &entity.User{
		ID:         1,
		BusinessID: 1,
		Email:      stringPtr("test@example.com"),
		FullName:   "Test User",
	}

	ctx := context.Background()

	testCases := []struct {
		name     string
		mock     func(m mocksForExecution)
		args     args
		expected expected
	}{
		{
			name: "positive: user successfully retrieved",
			mock: func(m mocksForExecution) {
				m.userRepo.On("Get", ctx, user.ID).Return(user, nil)
			},
			args: args{
				id: user.ID,
			},
			expected: expected{
				user: user,
			},
		},
		{
			name: "negative: user not found",
			mock: func(m mocksForExecution) {
				m.userRepo.On("Get", ctx, 999).Return(nil, repository.ErrNotFound)
			},
			args: args{
				id: 999,
			},
			expected: expected{
				err: fmt.Errorf("failed to get user: %w", repository.ErrNotFound),
			},
		},
		{
			name: "negative: repository error",
			mock: func(m mocksForExecution) {
				m.userRepo.On("Get", ctx, user.ID).Return(nil, fmt.Errorf("some error"))
			},
			args: args{
				id: user.ID,
			},
			expected: expected{
				err: fmt.Errorf("failed to get user: %w", fmt.Errorf("some error")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Init mocks
			userRepoMock := mocks.NewUserRepository(t)

			// Setup mocks
			tc.mock(mocksForExecution{
				userRepo: userRepoMock,
			})

			// Init service
			userService := services.NewUserService(&repository.Repositories{
				User: userRepoMock,
			})

			// Execute
			got, err := userService.Get(ctx, tc.args.id)

			// Assert
			if tc.expected.err != nil {
				assert.EqualError(t, err, tc.expected.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.user.ID, got.ID)
				assert.Equal(t, tc.expected.user.BusinessID, got.BusinessID)
				assert.Equal(t, tc.expected.user.Email, got.Email)
				assert.Equal(t, tc.expected.user.FullName, got.FullName)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
