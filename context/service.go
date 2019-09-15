package context

import "net/http"

//Represents a service instance.
type Serve func(log Log) (ServiceResult, error)

//Results of a service.
type ServiceResult struct {
	HttpHandler http.Handler
}
