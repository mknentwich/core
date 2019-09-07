package auth

import (
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/utils"
	"net/http"
)

//Logging function for this package.
var log context.Log

//Serve function for this package.
func Serve(logger context.Log) (context.ServiceResult, error) {
	log = logger
	mux := http.NewServeMux()
	mux.HandleFunc("/", utils.HttpImplement(log))
	mux.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {

	})
	return context.ServiceResult{HttpHandler: mux}, nil
}
