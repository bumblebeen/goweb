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
	"github.com/urfave/negroni"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/auth0/go-jwt-middleware"
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
	vars := mux.Vars(req);
	id := vars["id"];
	fmt.Println(id)
	fmt.Fprintf(res, "pong");
};

func bar (res http.ResponseWriter, req * http.Request) {
	vars := mux.Vars(req);
	id := vars["id"];
	fmt.Println(id)
	fmt.Fprintf(res, "bar");
};

func getSession() *mgo.Session{
	session, err := mgo.Dial("mongodb://localhost");
	if (err != nil) {
		panic(err);
	}
	return session
}

func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("MyMiddlware: before")
	next(rw, r)
	fmt.Println("MyMiddlware: after")
}

func MyMiddleware2(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Two: before")
	next(rw, r)
	fmt.Println("two: after")
}

func MyMiddleware3(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("-----------: before")
	next(rw, r)
	fmt.Println("-----------: after")
}

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	}
})

func main() {
	if (port == "") {
		port = ":8080"
	}
	fmt.Println("Starting server at port: ", port);

	logger := log.New(os.Stdout, "server: ", log.Lshortfile)

	router := mux.NewRouter();

	n := negroni.Classic()
	n.UseHandler(router)

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

	// API ROUTER
	subrouter := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)
	subrouter.HandleFunc("/", pong)
	subrouter.HandleFunc("/{id}", pong)
	subrouter.HandleFunc("/bar/{id}", bar)

	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(MyMiddleware),
		negroni.HandlerFunc(MyMiddleware2),
		negroni.Wrap(subrouter),
	))


	// SUB ROUTES using Common Middleware
	subrouter2 := mux.NewRouter().PathPrefix("/sub").Subrouter().StrictSlash(true)
	subrouter2.HandleFunc("/", pong)
	subrouter2.HandleFunc("/{id}", pong)
	subrouter2.HandleFunc("/bar/{id}", bar)

	common := negroni.New(
		negroni.HandlerFunc(MyMiddleware3),
	)

	router.PathPrefix("/sub").Handler(common.With(
		negroni.HandlerFunc(MyMiddleware),
		negroni.HandlerFunc(MyMiddleware2),
		negroni.Wrap(subrouter2),
	))

	// PANIC
	router.HandleFunc("/panic", func(res http.ResponseWriter, req *http.Request) {
		panic("PANIC!!!")
	});


	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	router.Handle("/patingin", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(myHandler),
	))

	s := &http.Server{
		Addr:           port,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
	//http.ListenAndServe(":8080", router);
}
