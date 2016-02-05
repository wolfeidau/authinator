package users

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"time"

	"github.com/wolfeidau/authinator/models"
)

var _ UserStore = &UserStoreLocal{}

// UserStoreLocal local user store for testing purposes
//
// NOTE: This will return passwords because pointers are used.
type UserStoreLocal struct {
	users map[string]*models.User
}

// NewUserStoreLocal create a new local user store
func NewUserStoreLocal() UserStore {
	return &UserStoreLocal{users: make(map[string]*models.User)}
}

// GetByID lookup a user by thier Identifier
func (usl *UserStoreLocal) GetByID(userID string) (*models.User, error) {
	usr, ok := usl.users[userID]

	if !ok {
		return nil, ErrUserNotFound
	}

	return usr, nil
}

// GetByLogin lookup a user by their login
func (usl *UserStoreLocal) GetByLogin(login string) (*models.User, error) {
	for _, v := range usl.users {
		if reflect.DeepEqual(v.Login, models.String(login)) {
			return v, nil
		}
	}
	return nil, ErrUserNotFound

}

// GetPasswordByLogin retrieve the users password for authentication
func (usl *UserStoreLocal) GetPasswordByLogin(login string) (string, error) {
	for _, v := range usl.users {
		if reflect.DeepEqual(v.Login, models.String(login)) {
			return models.StringValue(v.Password), nil
		}
	}
	return "", ErrUserNotFound
}

// Create create a new user in the system with the given information
func (usl *UserStoreLocal) Create(user *models.User) (*models.User, error) {

	var id string

	// check for unique login
	if usl.loginExists(models.StringValue(user.Login)) {
		return nil, ErrUserAlreadyExists
	}

	if user.ID == nil {
		id = newID()
		user.ID = models.String(id)
	} else {
		id = models.StringValue(user.ID)
	}

	usl.users[id] = user

	return user, nil
}

// Update the user, this is currently limited to changing the users name.
func (usl *UserStoreLocal) Update(user *models.User) error {
	cusr, ok := usl.users[models.StringValue(user.ID)]

	if !ok {
		return ErrUserNotFound
	}

	cusr.Name = user.Name

	return nil
}

// Delete delete the user by user ID
func (usl *UserStoreLocal) Delete(userID string) error {

	delete(usl.users, userID)

	return nil
}

// Exists Check if a user exists using the users login
func (usl *UserStoreLocal) Exists(login string) (bool, error) {
	return usl.loginExists(login), nil
}

func (usl *UserStoreLocal) loginExists(login string) bool {
	for _, v := range usl.users {
		if reflect.DeepEqual(v.Login, models.String(login)) {
			return true
		}
	}
	return false
}

func newID() string {
	h := sha1.New()

	c := []byte(time.Now().String())

	return fmt.Sprintf("%x", h.Sum(c))
}
