package core

import (
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/rest"
	"log"
)

//A map for all services which exist in that project.
var services = map[string]context.Serve{
	"api": rest.Serve,
	"db":  database.Serve}

//Calls the context to initialize everything.
func Serve() {
	err := context.Initialize(services)
	if err != nil {
		log.Fatalf("Fatal while start up application: %s", err.Error())
	}
}

//Calls the context to initialize everything with a custom host and port
func ServeCustomConfig(host string, sqLiteFile string) {
	err := context.InitializeCustomConfig(services, host, sqLiteFile)
	if err != nil {
		log.Fatalf("Fatal while start up application: %s", err.Error())
	}
}
