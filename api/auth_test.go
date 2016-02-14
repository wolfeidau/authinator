package api

import (
	"bytes"
	"testing"

	"github.com/wolfeidau/authinator/auth"
	"github.com/wolfeidau/authinator/store/users"
)

var userHash = "LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP"

func TestAuthenticateUser(t *testing.T) {

	certs, err := auth.GenerateTestCerts()
	if err != nil {
		t.Errorf("error generating test certs %v", err)
	}

	store := users.NewUserStoreLocal()

	store.Create(NewUser())

	ws := NewAuthResource(store, nil, certs)

	req := newFormRequest("POST", "http://api.his.com/users", bytes.NewBufferString("login=wolfeidau&password=Somewh3r3 there is a cow!"))

	recorder, resp := newResponse()

	ws.authenticateUser(req, resp)

	if recorder.Code != 200 {
		t.Errorf("expected 200 got %d %s", recorder.Code, recorder.Body.String())
	}

	if recorder.Header().Get("Authorization") == "" {
		t.Errorf("expected authorization header to exist")
	}

}
