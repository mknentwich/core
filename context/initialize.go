package context

import (
	"fmt"
	"log"
	"net/http"
)

//Initializes the context.
func initializeConfig() error {
	config, err := config()
	Conf = *config
	return err
}

//Initializes the custom context.
func initializeCustomConfig(host string, sqLiteFile string) error {
	config, err := customConfig(host, sqLiteFile)
	Conf = *config
	return err
}

//Initializes the whole context.
func Initialize(services map[string]Serve) error {
	err := initializeConfig()
	serveServices(services)
	if err != nil {
		return err
	}
	logger.Printf("listen on %s", Conf.Host)
	return http.ListenAndServe(Conf.Host, nil)
}

//Initializes the whole context with a custom configuration
func InitializeCustomConfig(services map[string]Serve, host string, sqLiteFile string) error {
	err := initializeCustomConfig(host, sqLiteFile)
	serveServices(services)
	if err != nil {
		return err
	}
	logger.Printf("listen on %s", Conf.Host)
	return http.ListenAndServe(Conf.Host, nil)
}

//Loops through all registered services, calls their Serve function and register their http handler.
func serveServices(services map[string]Serve) {
	for serviceId, serve := range services {
		result, err := serve(func(level LogLevel, format string, a ...interface{}) {
			logger.Printf("%s: %s: %s", level, serviceId, fmt.Sprintf(format, a...))
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		if result.HttpHandler != nil {
			http.Handle("/"+serviceId+"/", http.StripPrefix("/"+serviceId, result.HttpHandler))
		}
	}
}
