package api

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/gorilla/schema"
	"github.com/wolfeidau/authinator/auth"
	"github.com/wolfeidau/authinator/models"
	"github.com/wolfeidau/authinator/store/users"
	"github.com/wolfeidau/authinator/util"
)

var decoder = schema.NewDecoder()

// AuthResource user resource
type AuthResource struct {
	store      users.UserStore
	certs      *auth.Certs
	authFilter restful.FilterFunction
}

// NewAuthResource create a new user resource
func NewAuthResource(store users.UserStore, authFilter restful.FilterFunction, certs *auth.Certs) *AuthResource {
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

	creds := new(models.Authentication)
	err = decoder.Decode(creds, req.Request.PostForm)
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	phash, err := ar.store.GetPasswordByLogin(creds.Login)
	if err != nil {
		if err == users.ErrUserNotFound {
			resp.WriteHeaderAndEntity(http.StatusForbidden, errorMsg("Auth failed."))
			return
		}

		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	ok, err := util.CompareHashPassword(creds.Password, phash)
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	if !ok {
		resp.WriteHeaderAndEntity(http.StatusForbidden, errorMsg("Auth failed."))
		return
	}

	usr, err := ar.store.GetByLogin(creds.Login)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	tok, err := auth.GenerateClaim(ar.certs, usr)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	resp.AddHeader("Authorization", fmt.Sprintf("Bearer %s", tok))

	resp.WriteHeader(http.StatusOK)
}
