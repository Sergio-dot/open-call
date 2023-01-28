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

	query := `select id, email, password, created_at, updated_at from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
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

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set email = $1, updated_at = $2`

	_, err := m.DB.ExecContext(ctx, query, u.Email, time.Now())

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
	var id int
	var hashedPassword string

	// query the database
	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	// compares hashedPassword (stored in database) with plain text password (input by the user in the form)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
