package rest

import (
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/utils"
	"net/http"
)

//Logging function for this package.
var log context.Log

//Serve function for this package.
func Serve(logger context.Log) context.ServiceResult {
	log = logger
	mux := http.NewServeMux()
	mux.HandleFunc("/", utils.HttpImplement(log))
	return context.ServiceResult{HttpHandler: mux}
}
