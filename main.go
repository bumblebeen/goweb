package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/julienschmidt/httprouter"
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

	router := httprouter.New();
	uc := controllers.NewUserController(getSession());
	router.GET("/ping", func(res http.ResponseWriter, req * http.Request, p httprouter.Params){
		fmt.Fprintf(res, "pong");
	});
	router.GET("/user/:id", uc.GetUser);
	router.POST("/user", uc.CreateUser);
	router.POST("/user/login", uc.AuthenticateUser);
	router.DELETE("/user/:id", uc.RemoveUser);

	http.ListenAndServe(port, router);
}
