package controllers

import (
	"net/http"
	"github.com/bumblebeen/goweb/models"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"io"
	"log"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type UserController struct {
	session *mgo.Session;
}

func NewUserController(session *mgo.Session) *UserController {
	return &UserController{session}
}

func (uc UserController) GetUser (res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req);
	id := vars["id"];

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

	//res.Header().Set("Content-Type", "application/json");
	res.WriteHeader(200);
	fmt.Fprintf(res, "%s", uj)
}

func (uc UserController) AuthenticateUser (res http.ResponseWriter, req *http.Request) {
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

func (uc UserController) CreateUser (res http.ResponseWriter, req *http.Request) {
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

func (uc UserController) RemoveUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req);
	id := vars["id"];

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

func (uc UserController) GetTokenHandler (res http.ResponseWriter, r *http.Request){
	var mySigningKey = []byte("secret")
	/* Create the token */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	fmt.Println(token)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(tokenString)
	/* Finally, write the token to the browser window */
	res.Write([]byte(tokenString))
}

func (uc UserController) DecodeToken (res http.ResponseWriter, req *http.Request){
	var mySigningKey = []byte("secret")
	/* Create the token */
	vars := mux.Vars(req);
	tokenString := vars["token"];


	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return mySigningKey, nil
	})
	if (err != nil) {
		res.WriteHeader(http.StatusUnauthorized);
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("Type is %T\n", claims)
		fmt.Println(claims["foo"], claims["nbf"])
		res.WriteHeader(http.StatusOK);
	} else {
		res.WriteHeader(http.StatusUnauthorized);
	}
}
