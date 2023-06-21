package models

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type User struct {
	ID             int       `db:"id"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	HashedPassword []byte    `db:"password"`
	Created        time.Time `db:"created"`
}

type Users struct {
	Pool *pgxpool.Pool
}

func (m *Users) Insert(name, email, password string) error {
	return nil
}

func (m *Users) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *Users) Exists(id int) (bool, error) {
	return false, nil
}
