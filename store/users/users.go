package users

import (
	"errors"

	"github.com/wolfeidau/authinator/models"
)

var (
	ErrUserNotFound      = errors.New("User not found.")
	ErrUserAlreadyExists = errors.New("User already exists.")
)

// UserStore user store interface
type UserStore interface {
	GetByID(userID string) (*models.User, error)
	GetByLogin(login string) (*models.User, error)
	GetPasswordByLogin(login string) (string, error)
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) error
	Delete(userID string) error
	Exists(login string) (bool, error)
}
