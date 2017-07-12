package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/julienschmidt/httprouter"
	"github.com/bumblebeen/goweb/controllers"
)

var port = os.Getenv("PORT");

func main() {
	if (port == "") {
		port = ":8080"
	}
	fmt.Println("Starting server at port: ", port);

	router := httprouter.New();
	uc := controllers.NewUserController();
	router.GET("/ping", func(res http.ResponseWriter, req * http.Request, p httprouter.Params){
		fmt.Fprintf(res, "pong");
	});
	router.GET("/user/:id", uc.GetUser);
	router.POST("/user", uc.CreateUser);
	router.DELETE("/user/:id", uc.RemoveUser);

	http.ListenAndServe(port, router);
}
