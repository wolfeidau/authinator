package users

import (
	r "github.com/dancannon/gorethink"
	"github.com/wolfeidau/authinator/models"
)

var _ UserStore = &UserStoreRethinkDB{}

var (
	// DBName is the name of the RethinkDB database
	DBName = "authinator"
	// TableName is the name of users table in the RethinkDB database
	TableName = "users"
)

// UserStoreRethinkDB RethinkDB based user store
type UserStoreRethinkDB struct {
	session *r.Session
}

// GetByID retrieve a user from RethinkDB
func (us *UserStoreRethinkDB) GetByID(userID string) (*models.User, error) {

	res, err := r.DB(DBName).Table(TableName).Get(userID).Run(us.session)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	if res.IsNil() {
		return nil, ErrUserNotFound
	}

	usr := new(models.User)
	err = res.One(&usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// GetByLogin retrieve a user from RethinkDB filtering by their login
func (us *UserStoreRethinkDB) GetByLogin(login string) (*models.User, error) {

	res, err := r.DB(DBName).Table(TableName).Filter(map[string]interface{}{
		"login": login,
	}).Run(us.session)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	if res.IsNil() {
		return nil, ErrUserNotFound
	}

	usr := new(models.User)
	err = res.One(&usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// GetPasswordByLogin retrieve the password for a user using their login
func (us *UserStoreRethinkDB) GetPasswordByLogin(login string) (string, error) {
	res, err := r.DB(DBName).Table(TableName).Filter(map[string]interface{}{
		"login": login,
	}).Run(us.session)
	if err != nil {
		return "", err
	}

	defer res.Close()

	if res.IsNil() {
		return "", ErrUserNotFound
	}

	usr := new(models.User)
	err = res.One(&usr)
	if err != nil {
		return "", err
	}

	return models.StringValue(usr.Password), nil
}

// Create create the user in RethinkDB
func (us *UserStoreRethinkDB) Create(user *models.User) (*models.User, error) {

	resp, err := r.DB(DBName).Table(TableName).Insert(user).RunWrite(us.session)

	if err != nil {
		return nil, err
	}

	user.ID = models.String(resp.GeneratedKeys[0])

	return user, nil
}

// Update the user in RethinkDB
func (us *UserStoreRethinkDB) Update(user *models.User) error {

	userID := models.StringValue(user.ID)

	if user.Name == nil {
		return nil
	}

	_, err := r.DB(DBName).Table(TableName).Get(userID).Update(map[string]interface{}{
		"name": models.StringValue(user.Name),
	}).RunWrite(us.session)
	if err != nil {
		return err
	}

	return nil
}

// Delete delete the user from the RethinkDB database.
func (us *UserStoreRethinkDB) Delete(userID string) error {
	_, err := r.DB(DBName).Table(TableName).Get(userID).Delete().RunWrite(us.session)
	if err != nil {
		return err
	}

	return nil
}

// Exists check if the user exists in the RethinkDB database.
func (us *UserStoreRethinkDB) Exists(login string) (bool, error) {

	res, err := r.DB(DBName).Table(TableName).Filter(map[string]interface{}{
		"login": login,
	}).Run(us.session)
	if err != nil {
		return false, err
	}

	defer res.Close()

	if res.IsNil() {
		return false, nil
	}

	return true, nil
}
