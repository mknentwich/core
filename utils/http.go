package utils

import (
	"github.com/mknentwich/core/context"
	"net/http"
)

const methods = "POST, GET, OPTIONS, PUT, DELETE, HEAD"

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
	return Cors(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		handler(writer, request)
	})
}

func Cors(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", methods)
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		writer.Header().Set("Allowed", methods)
		if request.Method != http.MethodOptions {
			handler(writer, request)
		}
	}
}
