package controllers

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/bumblebeen/goweb/models"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"io"
	"log"
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
		res.WriteHeader(http.StatusNotFound)
		return
	}
	u:= models.User{};
	oid := bson.ObjectIdHex(id)
	if err := uc.session.DB("webapi").C("users").FindId(oid).One(&u); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	uj, _ := json.Marshal(u)

	res.Header().Set("Content-Type", "application/json");
	res.WriteHeader(200);
	fmt.Fprintf(res, "%s", uj)
}

func (uc UserController) AuthenticateUser (res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	u:= models.User{};
	un := req.PostFormValue("Name")
	p := req.PostFormValue("Password")

	if un == "" {
		res.WriteHeader(http.StatusBadRequest);
		fmt.Fprintf(res, "Bad Request")
		return
	}

	if err := uc.session.DB("webapi").C("users").Find(bson.M{"name": un}).One(&u); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		log.Println("Passwords do not match:", err)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	uj, _ := json.Marshal(u)

	res.Header().Set("Content-Type", "application/json");
	res.WriteHeader(200);
	fmt.Fprintf(res, "%s", uj)
}

func (uc UserController) CreateUser (res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		panic(err)
	}
	u := models.User{}
	if err := json.Unmarshal(body, &u); err != nil {
		panic(err)
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	u.Password = string(bs)
	// json.NewDecoder(req.Body).Decode(&u);

	u.Id = bson.NewObjectId()

	uc.session.DB("webapi").C("users").Insert(u);

	uj, _ := json.Marshal(u);

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
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
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusNoContent);
}