package rest

import (
	"encoding/json"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/utils"
	"net/http"
)

//Logging function for this package.
var log context.Log

//Serve function for this package.
func Serve(args context.ServiceArguments) (context.ServiceResult, error) {
	log = args.Log
	mux := http.NewServeMux()
	mux.HandleFunc("/", utils.HttpImplement(log))
	mux.HandleFunc("/categories", utils.Rest(get(QueryCategoriesWithChildrenAndScores)))
	return context.ServiceResult{HttpHandler: mux}, nil
}

//Encodes a structure as JSON and returns it
func get(query DataQuery) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			encoder := json.NewEncoder(writer)
			err := encoder.Encode(query())
			if err != nil {
				log(context.LOG_ERROR, "An error occurred on return a REST GET request: %s", err.Error())
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
