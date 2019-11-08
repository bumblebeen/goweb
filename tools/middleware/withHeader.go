package middleware

import (
	"net/http"
	"fmt"
)

func WithHeader(key, value string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Adding Header")
			w.Header().Add(key, value)
			h.ServeHTTP(w, r)
		});
	}
}
