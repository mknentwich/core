package context

import (
	"net/http"
)

//Represents a service instance.
type Serve func(args ServiceArguments) (ServiceResult, error)

//Arguments to a service.
type ServiceArguments struct {
	GeneratedDirectory string
	Log                Log
}

//Results of a service.
type ServiceResult struct {
	HttpHandler http.Handler
}
