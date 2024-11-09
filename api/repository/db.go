package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	SQLC *sqlc.Queries
	PGX  *pgxpool.Pool

	loc *time.Location
}

type Options struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DBName string

	MinConnections int
	MaxConnections int

	Timezone string
}

func (o *Options) connStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_min_conns=%d pool_max_conns=%d", o.Host, o.Port, o.User, o.Pass, o.DBName, o.MinConnections, o.MaxConnections)
}

const (
	defaultTimezone = "UTC"
)

func NewDB(opts Options) (*DB, error) {
	if opts.Timezone == "" {
		opts.Timezone = defaultTimezone
	}

	loc, err := time.LoadLocation(opts.Timezone)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	poolConfig, err := pgxpool.ParseConfig(opts.connStr())
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(opts.MaxConnections)
	poolConfig.MinConns = int32(opts.MinConnections)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	return &DB{
		SQLC: queries,
		PGX:  pool,
		loc:  loc,
	}, nil
}

func (db *DB) Close() {
	db.PGX.Close()
}

func (db *DB) Now() pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  time.Now().In(db.loc),
		Valid: true,
	}
}

func (db *DB) Begin() (pgx.Tx, error) {
	return db.PGX.Begin(context.Background())
}

var (
	ErrNotFound      = errors.New("not found in db")
	ErrAlreadyExists = errors.New("already exists")
)

func (db *DB) HandleBasicErrors(err error) error {
	if err == nil {
		return nil
	}
	e := err.Error()

	if e == "no rows in result set" {
		return ErrNotFound
	} else if strings.Contains(e, "duplicate key value violates unique constraint") {
		return ErrAlreadyExists
	}

	return err
}

func (db *DB) ValidText(t string) pgtype.Text {
	return pgtype.Text{
		String: t,
		Valid:  true,
	}
}
