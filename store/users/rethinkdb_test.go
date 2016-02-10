package users

import (
	"errors"
	"fmt"
	"testing"

	r "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
	"github.com/wolfeidau/authinator/models"
)

func CreateDB(session *r.Session) {
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

func DropDB(session *r.Session) {

	resp, err := r.DBDrop(DBName).RunWrite(session)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("%d DB dropped, %d tables dropped\n", resp.DBsDropped, resp.TablesDropped)
}

func TestGetUserRethinkDB(t *testing.T) {

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		t.Errorf("error connecting to rethinkdb: %s", err)
	}

	CreateDB(session)

	userID, err := createUser(session)
	if err != nil {
		t.Errorf("error creating user in rethinkdb: %s", err)
	}

	userStore := &UserStoreRethinkDB{session}

	usr, err := userStore.GetByID(userID)

	if err != nil {
		t.Errorf("error getting user in rethinkdb: %s", err)
	}

	assert.Equal(t, models.StringValue(usr.Email), "mark@wolfe.id.au")
}

func TestGetUserByLoginRethinkDB(t *testing.T) {

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		t.Errorf("error connecting to rethinkdb: %s", err)
	}

	CreateDB(session)

	userID, err := createUser(session)
	if err != nil {
		t.Errorf("error creating user in rethinkdb: %s", err)
	}

	userStore := &UserStoreRethinkDB{session}

	usr, err := userStore.GetByLogin("wolfeidau")

	if err != nil {
		t.Errorf("error getting user in rethinkdb: %s", err)
	}

	assert.Equal(t, "mark@wolfe.id.au", models.StringValue(usr.Email))
	assert.Equal(t, userID, models.StringValue(usr.ID))
}

func TestGetPasswordByLoginRethinkDB(t *testing.T) {

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		t.Errorf("error connecting to rethinkdb: %s", err)
	}

	CreateDB(session)

	_, err = createUser(session)
	if err != nil {
		t.Errorf("error creating user in rethinkdb: %s", err)
	}

	userStore := &UserStoreRethinkDB{session}

	pass, err := userStore.GetPasswordByLogin("wolfeidau")

	if err != nil {
		t.Errorf("error getting user in rethinkdb: %s", err)
	}

	assert.Equal(t, "LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP", pass)
}

func TestCreateUserRethinkDB(t *testing.T) {

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		t.Errorf("error connecting to rethinkdb: %s", err)
	}

	CreateDB(session)

	userStore := &UserStoreRethinkDB{session}

	usr, err := userStore.Create(&models.User{
		Email:    models.String("mark@wolfe.id.au"),
		Login:    models.String("wolfeidau"),
		Name:     models.String("Mark Wolfe"),
		Password: models.String("LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP"),
	})

	if err != nil {
		t.Errorf("error getting user in rethinkdb: %s", err)
	}

	assert.NotNil(t, usr.ID)
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
