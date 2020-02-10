package dav

import (
	"github.com/mknentwich/core/context"
	"golang.org/x/net/webdav"
	"net/http"
)

var log context.Log

func Serve(args context.ServiceArguments) (context.ServiceResult, error) {
	log = args.Log
	initializeTree()
	result := context.ServiceResult{HttpHandler: &webdav.Handler{
		Prefix:     "/",
		FileSystem: &PhantomFileSystem{},
		LockSystem: &PhantomLockSystem{},
		Logger: func(request *http.Request, err error) {
			if err != nil {
				log(context.LOG_ERROR, "an error occurred during a request: %s", err.Error())
			} else {
				log(context.LOG_INFO, "request completed without compilations")
			}
		},
	}}
	return result, nil
}
