package users

import (
	"errors"
	"fmt"
	"testing"

	r "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
	"github.com/wolfeidau/authinator/models"
)

func TestGetUserRethinkDB(t *testing.T) {

	_, userStore, userID, err := createUserStoreAndSession()

	if assert.NoError(t, err, "connecting to rethinkdb") {

		usr, err := userStore.GetByID(userID)
		if assert.NoError(t, err, "getting user from rethinkdb") {

			if assert.NotNil(t, usr) {
				assert.Equal(t, models.StringValue(usr.Email), "mark@wolfe.id.au")
			}
		}

		usr, err = userStore.GetByID("123")

		if assert.Error(t, err) {
			assert.Equal(t, err, ErrUserNotFound)
		}

	}
}

func TestGetUserByLoginRethinkDB(t *testing.T) {

	_, userStore, userID, err := createUserStoreAndSession()

	if assert.NoError(t, err, "connecting to rethinkdb") {

		usr, err := userStore.GetByLogin("wolfeidau")

		if assert.Nil(t, err) {
			assert.Equal(t, "mark@wolfe.id.au", models.StringValue(usr.Email))
			assert.Equal(t, userID, models.StringValue(usr.ID))
		}

		usr, err = userStore.GetByLogin("nothere")

		if assert.Error(t, err) {
			assert.Equal(t, err, ErrUserNotFound)
		}
	}
}

func TestGetPasswordByLoginRethinkDB(t *testing.T) {

	_, userStore, _, err := createUserStoreAndSession()

	if assert.NoError(t, err, "connecting to rethinkdb") {

		pass, err := userStore.GetPasswordByLogin("wolfeidau")

		if assert.NoError(t, err) {
			assert.Equal(t, "LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP", pass)
		}

		pass, err = userStore.GetPasswordByLogin("nothere")

		if assert.Error(t, err) {
			assert.Equal(t, err, ErrUserNotFound)
		}

	}
}

func TestCreateUserRethinkDB(t *testing.T) {

	_, userStore, _, err := createUserStoreAndSession()

	if assert.NoError(t, err, "connecting to rethinkdb") {

		usr, err := userStore.Create(&models.User{
			Email:    models.String("mark@wolfe.id.au"),
			Login:    models.String("wolfeidau"),
			Name:     models.String("Mark Wolfe"),
			Password: models.String("LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP"),
		})

		if assert.Nil(t, err) {
			assert.NotNil(t, usr.ID)
		}
	}
}

func TestUpdateUserRethinkDB(t *testing.T) {

	_, userStore, userID, err := createUserStoreAndSession()

	if assert.NoError(t, err, "connecting to rethinkdb") {

		err = userStore.Update(&models.User{
			ID:   models.String(userID),
			Name: models.String("Mark Wolfy"),
		})

		assert.NoError(t, err, "updating user in rethinkdb")

		err = userStore.Update(&models.User{
			Name: models.String("Mark Wolfy"),
		})

		if assert.Error(t, err) {
			assert.Equal(t, err, ErrUserNotFound)
		}
	}

}

func TestDeleteUserRethinkDB(t *testing.T) {

	_, userStore, userID, err := createUserStoreAndSession()

	if assert.NoError(t, err, "connecting to rethinkdb") {

		err = userStore.Delete(userID)
		assert.Nil(t, err, "deleting user in rethinkdb")
	}
}

func TestUserExistsRethinkDB(t *testing.T) {

	_, userStore, _, err := createUserStoreAndSession()

	if assert.NoError(t, err, "connecting to rethinkdb") {

		exists, err := userStore.Exists("wolfeidau")

		if assert.Nil(t, err, "checking if user exists in rethinkdb") {
			assert.True(t, exists)
		}

		exists, err = userStore.Exists("nothere")

		if assert.NoError(t, err, "checking if user exists in rethinkdb") {
			assert.False(t, exists)
		}
	}
}

func createUserDB(session *r.Session) {
	DBName = "authinator_test"

	resp, err := r.DBCreate(DBName).RunWrite(session)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%d DB created\n", resp.DBsCreated)

	r.DB(DBName).TableCreate(TableName).Exec(session)

	fmt.Printf("Table created\n")

	dresp, err := r.DB(DBName).Table(TableName).Delete().RunWrite(session)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%d rows deleted\n", dresp.Deleted)
}

func createUserStoreAndSession() (*r.Session, UserStore, string, error) {

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		return nil, nil, "", err
	}

	createUserDB(session)

	userID, err := createUser(session)
	if err != nil {
		return nil, nil, "", err
	}

	userStore := &UserStoreRethinkDB{session}

	return session, userStore, userID, nil
}

func createUser(session *r.Session) (string, error) {
	resp, err := r.DB(DBName).Table(TableName).Insert(map[string]interface{}{
		"email":    "mark@wolfe.id.au",
		"login":    "wolfeidau",
		"name":     "Mark Wolfe",
		"password": "LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP",
	}).RunWrite(session)

	if err != nil {
		return "", err
	}

	fmt.Printf("%d row inserted, %d key generated\n", resp.Inserted, len(resp.GeneratedKeys))

	if len(resp.GeneratedKeys) != 1 {
		return "", errors.New("Key not generated")
	}

	return resp.GeneratedKeys[0], nil
}
