package utils

import (
	"github.com/mknentwich/core/context"
	"net/http"
)

//Displays a warning on http and log side to implement this call.
func HttpImplement(log context.Log) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		msg := "please implement me"
		log(context.LOG_WARNING, msg)
		w.Write([]byte(msg))
	}
}

func Rest(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		handler(writer, request)
	}
}
