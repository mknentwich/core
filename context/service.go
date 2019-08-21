package context

import "net/http"

//Represents a service instance.
type Serve func(log Log) ServiceResult

//Results of a service.
type ServiceResult struct {
	HttpHandler http.Handler
}
