package dav

import (
	"github.com/mknentwich/core/context"
	"golang.org/x/net/webdav"
)

var log context.Log

func Serve(args context.ServiceArguments) (context.ServiceResult, error) {
	log = args.Log
	result := context.ServiceResult{HttpHandler: &webdav.Handler{
		Prefix:     "/",
		FileSystem: &PhantomFileSystem{},
		LockSystem: &PhantomLockSystem{},
	}}
	return result, nil
}
