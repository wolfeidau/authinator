package api

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/schema"
	"github.com/wolfeidau/authinator/models"
	"github.com/wolfeidau/authinator/store/users"
	"github.com/wolfeidau/authinator/util"
)

var decoder = schema.NewDecoder()

// AuthCerts used by JWT to sign tokens
type AuthCerts struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// AuthResource user resource
type AuthResource struct {
	store      users.UserStore
	certs      *AuthCerts
	authFilter restful.FilterFunction
}

// NewAuthResource create a new user resource
func NewAuthResource(store users.UserStore, authFilter restful.FilterFunction, certs *AuthCerts) *AuthResource {
	return &AuthResource{store, certs, authFilter}
}

// Register register the user resource with the rest container.
func (ar AuthResource) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Consumes(restful.MIME_JSON)

	ws.Path("/auth").
		Doc("Auth services").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	//
	ws.Route(ws.POST("/login").Consumes("application/x-www-form-urlencoded").To(ar.authenticateUser).Doc("Get the current user").Operation("authenicateUser"))

	container.Add(ws)
}

func (ar AuthResource) authenticateUser(req *restful.Request, resp *restful.Response) {

	err := req.Request.ParseForm()
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	auth := new(models.Authentication)
	err = decoder.Decode(auth, req.Request.PostForm)
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	phash, err := ar.store.GetPasswordByLogin(auth.Login)
	if err != nil {
		if err == users.ErrUserNotFound {
			resp.WriteHeaderAndEntity(http.StatusForbidden, errorMsg("Auth failed."))
			return
		}

		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	ok, err := util.CompareHashPassword(auth.Password, phash)
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	if !ok {
		resp.WriteHeaderAndEntity(http.StatusForbidden, errorMsg("Auth failed."))
		return
	}

	usr, err := ar.store.GetByLogin(auth.Login)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	// generate a token
	var claims = jws.Claims{
		"email": models.StringValue(usr.Email),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	j := jws.NewJWT(claims, crypto.SigningMethodRS512)
	b, err := j.Serialize(ar.certs.PrivateKey)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	resp.AddHeader("Authorization", fmt.Sprintf("Bearer %s", string(b)))

	resp.WriteHeader(http.StatusOK)
}
