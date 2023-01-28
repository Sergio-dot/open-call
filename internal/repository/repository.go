package repository

import "github.com/Sergio-dot/open-call/internal/models"

type DatabaseRepo interface {
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
}
