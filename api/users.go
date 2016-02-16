package api

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/wolfeidau/authinator/models"
	"github.com/wolfeidau/authinator/store/users"
	"github.com/wolfeidau/authinator/util"
	"github.com/wolfeidau/authinator/validation"
)

// UserResource user resource
type UserResource struct {
	store      users.UserStore
	authFilter restful.FilterFunction
}

// NewUserResource create a new user resource
func NewUserResource(store users.UserStore, authFilter restful.FilterFunction) *UserResource {
	return &UserResource{store, authFilter}
}

// Register register the user resource with the rest container.
func (ur UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Consumes(restful.MIME_JSON)

	ws.Path("/users").
		Doc("User services").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").Filter(ur.authFilter).To(ur.getUser).
		Doc("Get the current user").
		Operation("getUser").Reads(models.User{}))

	ws.Route(ws.PUT("/").Filter(ur.authFilter).To(ur.updateUser).
		Doc("Update your user information").
		Operation("updateUser").Writes(models.User{}))

	ws.Route(ws.POST("/").To(ur.createUser).
		Doc("Register a new user").
		Operation("createUser").Writes(models.User{}))

	ws.Route(ws.PUT("/password").Filter(ur.authFilter).To(ur.updatePassword).
		Doc("Update the current users password").
		Operation("updatePassword").Writes(models.User{}))

	container.Add(ws)
}

func (ur UserResource) getUser(req *restful.Request, resp *restful.Response) {

	userid, ok := req.Attribute("user_id").(string)

	if !ok {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	usr, err := ur.store.GetByID(userid)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusNotFound, errorMsg("User not found."))
		return
	}

	resp.WriteEntity(usr)
}

func (ur UserResource) updateUser(req *restful.Request, resp *restful.Response) {

	userid, ok := req.Attribute("user_id").(string)

	if !ok {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	cusr, err := ur.store.GetByID(userid)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusNotFound, errorMsg("User not found."))
		return
	}

	usr := new(models.User)
	err = req.ReadEntity(usr)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	// ensure the userid is trusted
	usr.ID = models.String(userid)

	allErrs := validation.ValidateUserUpdate(usr, cusr)

	if len(allErrs) != 0 {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, validationErrors("validation failed", allErrs))
		return
	}

	err = ur.store.Update(usr)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	usr.Password = nil

	resp.WriteEntity(usr)
}

func (ur UserResource) createUser(req *restful.Request, resp *restful.Response) {

	usr := new(models.User)
	err := req.ReadEntity(usr)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	allErrs := validation.ValidateUserRegister(usr)

	if len(allErrs) != 0 {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, validationErrors("validation failed", allErrs))
		return
	}

	exists, err := ur.store.Exists(models.StringValue(usr.Login))

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	if exists {
		resp.WriteHeaderAndEntity(http.StatusConflict, errorMsg("User already exists."))
		return
	}

	// hash the password
	pass := models.StringValue(usr.Password)

	pass, err = util.HashPassword(pass)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	usr.Password = models.String(pass)

	nusr, err := ur.store.Create(usr)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	nusr.Password = nil

	resp.WriteHeaderAndEntity(http.StatusCreated, nusr)
}

func (ur UserResource) updatePassword(req *restful.Request, resp *restful.Response) {

	userid, ok := req.Attribute("user_id").(string)

	if !ok {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	cusr, err := ur.store.GetByID(userid)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusNotFound, errorMsg("User not found."))
		return
	}

	data := make(map[string]string)
	err = req.ReadEntity(&data)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	password, ok := data["password"]

	if !ok {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, errorMsg("bad request missing password"))
		return
	}

	// hash the password
	pass, err := util.HashPassword(password)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	cusr.Password = models.String(pass)

	err = ur.store.Update(cusr)

	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, errorMsg("Server error."))
		return
	}

	resp.WriteHeader(http.StatusOK)
}
