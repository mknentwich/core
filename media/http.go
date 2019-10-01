package media

import (
	"github.com/mknentwich/core/context"
	"net/http"
	"strconv"
	"strings"
)

var log context.Log

const (
	audio = "audio"
	pdf   = "pdf"
)

func Serve(args context.ServiceArguments) context.ServiceResult {
	log = args.Log
	mux := http.NewServeMux()
	mux.HandleFunc("/score/", httpScore)
	return context.ServiceResult{
		HttpHandler: mux}
}

func httpScore(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	_, err := strconv.Atoi(parts[len(parts)-2])
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	mediaType := parts[len(parts)-1]
	if mediaType != audio && mediaType != pdf {
		rw.WriteHeader(http.StatusNotFound)
		log(context.LOG_INFO, "someone tried to access unknown media type: %s", mediaType)
	}
}
