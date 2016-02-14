package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/wolfeidau/authinator/models"
	"github.com/wolfeidau/authinator/store/users"
)

const (
	newUserJSON        = `{"login":"wolfeidau","email":"mark@wolfe.id.au","name":"Mark Wolfe","password":"Somewh3r3 there is a cow!"}`
	updateUserJSON     = `{"login":"wolfeidau","email":"mark@wolfe.id.au","name":"Mark Wolf"}`
	updatePasswordJSON = `{"password":"Somewh3r3 there is a cow!"}`
)

func TestGetUser(t *testing.T) {

	_, ws := setupResourceAndStore()

	req := newRequest("PUT", "http://api.his.com/users", bytes.NewBufferString(updateUserJSON))

	// add a mock Attribute which should be added by the auth filter
	req.SetAttribute("user_id", "123")

	recorder, resp := newResponse()

	ws.updateUser(req, resp)

	if recorder.Code != 200 {
		t.Errorf("expected 200 got %d %s", recorder.Code, recorder.Body.String())
	}

}

func TestCreateUser(t *testing.T) {

	ws := NewUserResource(users.NewUserStoreLocal(), nil)

	req := newRequest("POST", "http://api.his.com/users", bytes.NewBufferString(newUserJSON))
	recorder, resp := newResponse()

	ws.createUser(req, resp)

	if recorder.Code != 201 {
		t.Errorf("expected 201 got %d %s", recorder.Code, recorder.Body.String())
	}
}

func TestUpdateUser(t *testing.T) {

	_, ws := setupResourceAndStore()

	req := newRequest("PUT", "http://api.his.com/users", bytes.NewBufferString(updateUserJSON))

	// add a mock Attribute which should be added by the auth filter
	req.SetAttribute("user_id", "123")

	recorder, resp := newResponse()

	ws.updateUser(req, resp)

	if recorder.Code != 200 {
		t.Errorf("expected 200 got %d %s", recorder.Code, recorder.Body.String())
	}
}

func TestUpdatePassword(t *testing.T) {

	_, ws := setupResourceAndStore()

	req := newRequest("PUT", "http://api.his.com/users/password", bytes.NewBufferString(updatePasswordJSON))

	// add a mock Attribute which should be added by the auth filter
	req.SetAttribute("user_id", "123")

	recorder, resp := newResponse()

	ws.updatePassword(req, resp)

	if recorder.Code != 200 {
		t.Errorf("expected 200 got %d %s", recorder.Code, recorder.Body.String())
	}
}

func setupResourceAndStore() (users.UserStore, *UserResource) {
	store := users.NewUserStoreLocal()

	store.Create(NewUser())

	ws := NewUserResource(store, nil)

	return store, ws
}

func newRequest(method, urlStr string, body io.Reader) *restful.Request {
	httpReq, _ := http.NewRequest(method, urlStr, body)
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	return restful.NewRequest(httpReq)
}

func newFormRequest(method, urlStr string, body io.Reader) *restful.Request {
	httpReq, _ := http.NewRequest(method, urlStr, body)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	return restful.NewRequest(httpReq)
}

func newResponse() (*httptest.ResponseRecorder, *restful.Response) {
	recorder := httptest.NewRecorder()
	resp := restful.NewResponse(recorder)
	resp.SetRequestAccepts(restful.MIME_JSON)
	return recorder, resp
}

func NewUser() *models.User {
	return &models.User{
		ID:       models.String("123"),
		Login:    models.String("wolfeidau"),
		Email:    models.String("mark@wolfe.id.au"),
		Name:     models.String("Mark Wolfe"),
		Password: models.String("LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP"), //Somewh3r3 there is a cow!
	}
}
