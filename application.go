package core

import (
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/dav"
	"github.com/mknentwich/core/media"
	"github.com/mknentwich/core/pdf"
	"github.com/mknentwich/core/rest"
	"github.com/mknentwich/core/template"
	"log"
)

//A map for all services which exist in that project.
var services = map[string]context.Serve{
	"api":      rest.Serve,
	"dav":      dav.Serve,
	"db":       database.Serve,
	"media":    media.Serve,
	"pdf":      pdf.Serve,
	"template": template.Serve}

//Calls the context to initialize everything.
func Serve() {
	err := context.Initialize(services)
	if err != nil {
		log.Fatalf("Fatal while start up application: %s", err.Error())
	}
}
