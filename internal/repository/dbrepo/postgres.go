package dbrepo

import (
	"context"
	"errors"
	"github.com/Sergio-dot/open-call/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// GetUserByID returns a User by its ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, username, email, password, created_at, updated_at from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// CreateUser creates a new record in the database representing a user
func (m *postgresDBRepo) CreateUser(u models.User) error {
	// create a temporary context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query the database
	query := `insert into users (email, username, password, created_at, updated_at) values ($1, $2, $3, $4, $5)`

	// generate hashed password from testPassword (input by the user)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 12)

	// exec query with context
	_, err := m.DB.ExecContext(ctx, query, u.Email, u.Username, hashedPassword, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set email = $1, username = $2, updated_at = $3 where id = $4`

	_, err := m.DB.ExecContext(ctx, query, u.Email, u.Username, time.Now(), u.ID)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate perform user authentication
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	// create a temporary context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// variables used to store information from database
	var (
		id             int
		username       string
		hashedPassword string
	)

	// query the database
	row := m.DB.QueryRowContext(ctx, "select id, username, password from users where email = $1", email)
	err := row.Scan(&id, &username, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	// compares hashedPassword (stored in database) with password (input by the user in the form) hashed
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
