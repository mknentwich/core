package auth

import (
	"encoding/json"
	"fmt"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"net/http"
	"strings"
)

//Logging function for this package.
var log context.Log

//Serve function for this package.
func Serve(args context.ServiceArguments) (context.ServiceResult, error) {
	log = args.Log
	mux := http.NewServeMux()
	mux.HandleFunc("/login", httpLogin)
	mux.HandleFunc("/password", httpPassword)
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

func httpPassword(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	requester, _, err := Check(getJwt(r))
	if err != nil {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err = decoder.Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if user.Email != requester.Email && requester.Admin {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	err = saveUser(&user)
	if err != nil {
		log(context.LOG_ERROR, "error occured while creating/updating user: %s", err.Error())
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

//Wrapper for HandlerFuncs to verify if user is logged in and if admin is necessary or not.
func authenticated(handlerFunc http.HandlerFunc, admin bool) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user, _, err := Check(getJwt(r))
		if !context.Conf.Authentication || (err == nil && user != nil && (user.Admin || !admin)) {
			handlerFunc(rw, r)
		} else {
			rw.WriteHeader(http.StatusForbidden)
		}
	}
}

//Checks if user is logged in. If not, the handlerFunc won't be called.
func Auth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return authenticated(handlerFunc, false)
}

//Checks if user is admin. If not, the handlerFunc won't be called.
func Admin(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return authenticated(handlerFunc, true)
}
