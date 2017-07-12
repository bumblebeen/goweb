package controllers

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/bumblebeen/goweb/models"
	"encoding/json"
	"fmt"
)

type UserController struct {}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc UserController) GetUser (res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	u:= models.User{
		Name: "Marvin Arcilla",
		Gender: "Male",
		Age: 24,
		Id: ps.ByName("id"),
	}

	uj, _ := json.Marshal(u)

	res.Header().Set("Content-Type", "application/json");
	res.WriteHeader(200);
	fmt.Fprintf(res, "%s", uj)
}

func (uc UserController) CreateUser (res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	u := models.User{}

	json.NewDecoder(req.Body).Decode(&u);

	u.Id = "foo"

	uj, _ := json.Marshal(u);

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(201)
	fmt.Fprintf(res, "%s", uj)
}

func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: only write status for now
	w.WriteHeader(200)
}