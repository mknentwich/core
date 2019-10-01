package media

import (
	"github.com/mknentwich/core/auth"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/utils"
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
	scoreId, err := strconv.Atoi(parts[len(parts)-2])
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	mediaType := parts[len(parts)-1]
	if mediaType != audio && mediaType != pdf {
		rw.WriteHeader(http.StatusNotFound)
		log(context.LOG_INFO, "someone tried to access unknown media type: %s", mediaType)
	}
	switch r.Method {
	case http.MethodDelete:
		auth.Auth(mediaDelete(scoreId, mediaType))(rw, r)
	case http.MethodGet:
		mediaGet(scoreId, mediaType)(rw, r)
	case http.MethodPost:
		auth.Auth(mediaPost(scoreId, mediaType))(rw, r)
	case http.MethodPut:
		auth.Auth(mediaPut(scoreId, mediaType))(rw, r)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func mediaDelete(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, rr *http.Request) {
		utils.HttpImplement(log)
	}
}

func mediaGet(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, rr *http.Request) {
		utils.HttpImplement(log)
	}
}

func mediaPost(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, rr *http.Request) {
		utils.HttpImplement(log)
	}
}

func mediaPut(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, rr *http.Request) {
		utils.HttpImplement(log)
	}
}