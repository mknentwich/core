package context

import "net/http"

//Represents a service instance.
type Service interface {
	Serve(config Configuration, log Log)
}

//Results of a service.
type ServiceResult struct {
	HttpMux *http.ServeMux
}
