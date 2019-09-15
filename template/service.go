package template

import (
	"github.com/mknentwich/core/auth"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/utils"
	"net/http"
)

//Logging function for this package.
var log context.Log

//Output directory for generated files.
var outDir string

//Serve function for this package.
func Serve(args context.ServiceArguments) (context.ServiceResult, error) {
	log = args.Log
	outDir = args.GeneratedDirectory
	mux := http.NewServeMux()
	mux.HandleFunc("/generate", auth.Admin(httpGenerate))
	return context.ServiceResult{HttpHandler: mux}, nil
}

//HTTP call to generate markup manually.
func httpGenerate(rw http.ResponseWriter, r *http.Request) {
	log(context.LOG_INFO, "manual template generation was issued")
	utils.HttpImplement(log)(rw, r)
}
