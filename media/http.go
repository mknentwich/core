package media

import (
	"github.com/mknentwich/core/auth"
	"github.com/mknentwich/core/context"
	"net/http"
	"strconv"
	"strings"
)

//log function for this package
var log context.Log

//directory where all media will be stored
var outDir string

//media types
const (
	audio = "audio"
	pdf   = "pdf"
)

//slice of media types
var mediaTypes = []string{audio, pdf}

//instantiates this package
func Serve(args context.ServiceArguments) (context.ServiceResult, error) {
	log = args.Log
	outDir = args.GeneratedDirectory
	err := createDirs()
	mux := http.NewServeMux()
	mux.HandleFunc("/score/", httpScore)
	return context.ServiceResult{
		HttpHandler: mux}, err
}

//root handler for scores
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

//returns a http handler which deletes a media from the filesystem
func mediaDelete(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !scoreExist(scoreId, mediaType) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		writeInternalError(removeMedia(scoreId, mediaType, rw), rw)
	}
}

//returns a http handler which reads a media from the filesystem and returns it in the body
func mediaGet(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, rr *http.Request) {
		contentType(mediaType, rw)
		if !scoreExist(scoreId, mediaType) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		writeInternalError(readMediaFromDiskTo(scoreId, mediaType, rw), rw)
	}
}

//returns a http handler which reads the media from the body and stores it (only for creation)
func mediaPost(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		writeInternalError(saveMediaToDisk(scoreId, mediaType, r.Body), rw)
	}
}

//returns a http handler which reads the media from the body and stores it (only for overriding)
func mediaPut(scoreId int, mediaType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !scoreExist(scoreId, mediaType) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		writeInternalError(saveMediaToDisk(scoreId, mediaType, r.Body), rw)
	}
}

//sets the correct Content-Type to the response
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
