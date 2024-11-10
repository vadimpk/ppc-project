package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/repository/db/sqlc"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --dir . --name UserRepository --output ./mocks
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	CreateBusinessAdmin(ctx context.Context, businessName string, user *entity.User) error
	Get(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, businessID int, email string) (*entity.User, error)
	GetByPhone(ctx context.Context, businessID int, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	UpdatePassword(ctx context.Context, id int, passwordHash string) error
}

type userRepository struct {
	db *DB
}

func NewUserRepository(db *DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	var email, phone pgtype.Text
	if user.Email != nil {
		email = pgtype.Text{String: *user.Email, Valid: true}
	}
	if user.Phone != nil {
		phone = pgtype.Text{String: *user.Phone, Valid: true}
	}

	params := sqlc.CreateUserParams{
		BusinessID:   pgtype.Int4{Int32: int32(user.BusinessID), Valid: true},
		Email:        email,
		Phone:        phone,
		FullName:     user.FullName,
		PasswordHash: r.db.ValidText(user.PasswordHash),
		Role:         user.Role,
	}

	dbUser, err := r.db.SQLC.CreateUser(ctx, params)
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	user.ID = int(dbUser.ID)
	user.CreatedAt = dbUser.CreatedAt.Time
	return nil
}

func (r *userRepository) CreateBusinessAdmin(ctx context.Context, businessName string, user *entity.User) error {
	var email, phone pgtype.Text
	if user.Email != nil {
		email = pgtype.Text{String: *user.Email, Valid: true}
	}
	if user.Phone != nil {
		phone = pgtype.Text{String: *user.Phone, Valid: true}
	}

	params := sqlc.CreateBusinessAdminParams{
		Name:         businessName,
		Email:        email,
		Phone:        phone,
		FullName:     user.FullName,
		PasswordHash: r.db.ValidText(user.PasswordHash),
	}

	dbUser, err := r.db.SQLC.CreateBusinessAdmin(ctx, params)
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	user.ID = int(dbUser.ID)
	user.BusinessID = int(dbUser.BusinessID.Int32)
	user.CreatedAt = dbUser.CreatedAt.Time
	return nil
}

func (r *userRepository) Get(ctx context.Context, id int) (*entity.User, error) {
	dbUser, err := r.db.SQLC.GetUserByID(ctx, int32(id))
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	return convertDBUserToEntity(dbUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, businessID int, email string) (*entity.User, error) {
	params := sqlc.GetUserByEmailParams{
		BusinessID: pgtype.Int4{
			Int32: int32(businessID),
			Valid: true,
		},
		Email: pgtype.Text{String: email, Valid: true},
	}

	dbUser, err := r.db.SQLC.GetUserByEmail(ctx, params)
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	return convertDBUserToEntity(dbUser), nil
}

func (r *userRepository) GetByPhone(ctx context.Context, businessID int, phone string) (*entity.User, error) {
	params := sqlc.GetUserByPhoneParams{
		BusinessID: pgtype.Int4{
			Int32: int32(businessID),
			Valid: true,
		},
		Phone: pgtype.Text{String: phone, Valid: true},
	}

	dbUser, err := r.db.SQLC.GetUserByPhone(ctx, params)
	if err != nil {
		return nil, r.db.HandleBasicErrors(err)
	}

	return convertDBUserToEntity(dbUser), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	var email, phone pgtype.Text
	if user.Email != nil {
		email = pgtype.Text{String: *user.Email, Valid: true}
	}
	if user.Phone != nil {
		phone = pgtype.Text{String: *user.Phone, Valid: true}
	}

	params := sqlc.UpdateUserParams{
		ID:       int32(user.ID),
		Email:    email,
		Phone:    phone,
		FullName: user.FullName,
	}

	dbUser, err := r.db.SQLC.UpdateUser(ctx, params)
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	user.CreatedAt = dbUser.CreatedAt.Time
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	params := sqlc.UpdateUserPasswordParams{
		ID:           int32(id),
		PasswordHash: r.db.ValidText(passwordHash),
	}

	err := r.db.SQLC.UpdateUserPassword(ctx, params)
	if err != nil {
		return r.db.HandleBasicErrors(err)
	}

	return nil
}

func convertDBUserToEntity(dbUser sqlc.User) *entity.User {
	user := &entity.User{
		ID:           int(dbUser.ID),
		BusinessID:   int(dbUser.BusinessID.Int32),
		FullName:     dbUser.FullName,
		PasswordHash: dbUser.PasswordHash.String,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt.Time,
	}

	if dbUser.Email.Valid {
		email := dbUser.Email.String
		user.Email = &email
	}
	if dbUser.Phone.Valid {
		phone := dbUser.Phone.String
		user.Phone = &phone
	}

	return user
}
