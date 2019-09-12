package auth

import (
	"encoding/json"
	"fmt"
	"github.com/mknentwich/core/context"
	"net/http"
	"strings"
)

//Logging function for this package.
var log context.Log

//Serve function for this package.
func Serve(logger context.Log) (context.ServiceResult, error) {
	log = logger
	mux := http.NewServeMux()
	mux.HandleFunc("/login", httpLogin)
	mux.HandleFunc("/refresh", httpRefresh)
	mux.HandleFunc("/self", httpSelf)
	return context.ServiceResult{HttpHandler: mux}, nil
}

func httpLogin(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	credentials := Credentials{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	jwt, err := login(&credentials)
	if err != nil {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	_, err = rw.Write([]byte(jwt))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func httpRefresh(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	jwt := getJwt(r)
	user, _, err := Check(jwt)
	if err != nil {
		fmt.Println(err.Error())
		if _, ok := err.(*jwtFormatError); ok {
			rw.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			rw.WriteHeader(http.StatusForbidden)
		}
		return
	}
	_, err = rw.Write([]byte(generateToken(user)))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func httpSelf(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	jwt := getJwt(r)
	user, _, err := Check(jwt)
	if err != nil {
		fmt.Println(err.Error())
		if _, ok := err.(*jwtFormatError); ok {
			rw.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			rw.WriteHeader(http.StatusForbidden)
		}
		return
	}
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

//Get JWT from HTTP request.
func getJwt(r *http.Request) string {
	return strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
}
