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
func Serve(logger context.Log) (context.ServiceResult, error) {
	log = logger
	mux := http.NewServeMux()
	mux.HandleFunc("/", utils.HttpImplement(log))
	mux.HandleFunc("/categories", utils.Rest(get(QueryCategoriesWithChildrenAndScores)))
	mux.HandleFunc("/scores", utils.Rest(get(QueryScoresFlat)))
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

//handles http handlers in a sequence.
func multi(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler(rw, r)
		}
	}
}
