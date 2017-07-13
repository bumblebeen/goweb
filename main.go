package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/bumblebeen/goweb/controllers"
	"gopkg.in/mgo.v2"
)

var port = os.Getenv("PORT");

func getSession() *mgo.Session{
	session, err := mgo.Dial("mongodb://localhost");
	if (err != nil) {
		panic(err);
	}
	return session
}

func main() {
	if (port == "") {
		port = ":8080"
	}
	fmt.Println("Starting server at port: ", port);

	router := mux.NewRouter();
	uc := controllers.NewUserController(getSession());
	router.HandleFunc("/ping", func(res http.ResponseWriter, req * http.Request){
		fmt.Fprintf(res, "pong");
	});
	router.HandleFunc("/user/{id}", uc.GetUser).Methods("GET");
	router.HandleFunc("/token", uc.GetTokenHandler).Methods("GET");
	router.HandleFunc("/token/validate/{token}", uc.DecodeToken).Methods("GET");
	router.HandleFunc("/user", uc.CreateUser).Methods("POST");
	router.HandleFunc("/user/login", uc.AuthenticateUser).Methods("POST");
	router.HandleFunc("/user/{id}", uc.RemoveUser).Methods("DELETE");

	http.ListenAndServe(port, router);
}
