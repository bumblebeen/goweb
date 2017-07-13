package middleware

import (
	"net/http"
	"fmt"
)

func SampleMw () http.Handler {
	return http.HandlerFunc(func (res http.ResponseWriter, req *http.Request){
		fmt.Println("wooii")
	});
}
