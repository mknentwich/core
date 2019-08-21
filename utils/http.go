package utils

import (
	"github.com/mknentwich/core/context"
	"net/http"
)

//Displays a warning on http and log side to implement this call.
func HttpImplement(log context.Log) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := "please implement me"
		log(context.LOG_WARNING, msg)
		w.Write([]byte(msg))
	}
}
