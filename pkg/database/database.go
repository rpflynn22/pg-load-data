package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// GetDB returns a DB connection initialized to talk to a local postgres.
func GetDB(ctx context.Context) (*sqlx.DB, error) {
	cfg, err := pgx.ParseConfig("postgres://postgres:@localhost:5432/postgres")
	if err != nil {
		return nil, err
	}
	sqlDB := stdlib.OpenDB(*cfg)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	db := sqlx.NewDb(sqlDB, "pgx")
	db.SetMaxOpenConns(99)
	db.SetMaxIdleConns(99)
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
