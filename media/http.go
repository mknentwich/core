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

var outDir string

const (
	audio = "audio"
	pdf   = "pdf"
)

var mediaTypes = []string{audio, pdf}

func Serve(args context.ServiceArguments) (context.ServiceResult, error) {
	log = args.Log
	outDir = args.GeneratedDirectory
	err := createDirs()
	mux := http.NewServeMux()
	mux.HandleFunc("/score/", httpScore)
	return context.ServiceResult{
		HttpHandler: mux}, err
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
		if !scoreExist(scoreId) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		utils.HttpImplement(log)
	}
}

func mediaGet(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, rr *http.Request) {
		contentType(mediaType, rw)
		if !scoreExist(scoreId) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		writeInternalError(readMediaFromDiskTo(scoreId, mediaType, rw), rw)
	}
}

func mediaPost(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		writeInternalError(saveMediaToDisk(scoreId, mediaType, r.Body), rw)
	}
}

func mediaPut(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !scoreExist(scoreId) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		writeInternalError(saveMediaToDisk(scoreId, mediaType, r.Body), rw)
	}
}

func contentType(mediaType string, rw http.ResponseWriter) {
	ct := ""
	switch mediaType {
	case pdf:
		ct = "application/pdf"
	case audio:
		ct = "audio/mpeg"
	}
	rw.Header().Set("Content-Type", ct)
}
