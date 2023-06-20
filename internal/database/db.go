package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

// connectDB tries to open a connection and returns it.
func Connect(dsn string) (*DB, error) {
	db, err := pgxpool.New(context.TODO(), dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.TODO())
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
