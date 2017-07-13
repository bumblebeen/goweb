package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/bumblebeen/goweb/controllers"
	"gopkg.in/mgo.v2"
	"log"
	"github.com/bumblebeen/goweb/tools/middleware"
)

var port = os.Getenv("PORT");

func notifyMw(logger *log.Logger) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("before")
			defer logger.Println("after")
			h.ServeHTTP(w, r)
		});
	}
}

func pong (res http.ResponseWriter, req * http.Request) {
	fmt.Fprintf(res, "pong");
};


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

	logger := log.New(os.Stdout, "server: ", log.Lshortfile)

	router := mux.NewRouter();
	uc := controllers.NewUserController(getSession());

	router.HandleFunc("/ping", pong);
	router.HandleFunc("/token", uc.GetTokenHandler).Methods("GET");
	router.HandleFunc("/token/validate/{token}", uc.DecodeToken).Methods("GET");
	router.HandleFunc("/user", uc.CreateUser).Methods("POST");
	router.HandleFunc("/user/login", uc.AuthenticateUser).Methods("POST");
	router.HandleFunc("/user/{id}", uc.RemoveUser).Methods("DELETE");

	router.Handle("/user/{id}", middleware.HandleMiddleWares(
		http.HandlerFunc(uc.GetUser),
		middleware.WithHeader("Content-Type", "application/json"),
	)).Methods("GET");

	router.Handle("/middlewares", middleware.HandleMiddleWares(
		middleware.SampleMw(),
		notifyMw(logger),
		middleware.WithHeader("Content-Type", "application/json"),
	)).Methods("GET")

	http.ListenAndServe(port, router);
}
