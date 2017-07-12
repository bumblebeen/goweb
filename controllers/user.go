package controllers

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/bumblebeen/goweb/models"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session;
}

func NewUserController(session *mgo.Session) *UserController {
	return &UserController{session}
}

func (uc UserController) GetUser (res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if !bson.IsObjectIdHex(id) {
		res.WriteHeader(404)
		return
	}
	u:= models.User{};
	oid := bson.ObjectIdHex(id)
	if err := uc.session.DB("webapi").C("users").FindId(oid).One(&u); err != nil {
		res.WriteHeader(404)
		return
	}

	uj, _ := json.Marshal(u)

	res.Header().Set("Content-Type", "application/json");
	res.WriteHeader(200);
	fmt.Fprintf(res, "%s", uj)
}

func (uc UserController) CreateUser (res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	u := models.User{}

	json.NewDecoder(req.Body).Decode(&u);

	u.Id = bson.NewObjectId()

	uc.session.DB("webapi").C("users").Insert(u);

	uj, _ := json.Marshal(u);

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(201)
	fmt.Fprintf(res, "%s", uj)
}

func (uc UserController) RemoveUser(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if !bson.IsObjectIdHex(id) {
		res.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(id)
	if err := uc.session.DB("webapi").C("users").RemoveId(oid); err != nil {
		res.WriteHeader(404)
		return
	}

	res.WriteHeader(204);
}