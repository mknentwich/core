package utils

import (
	"path"
	"text/template"
)

const tmplDir = "templates"

func CreateTemplate(name string) (*template.Template, error) {
	return template.New(name).ParseFiles(path.Join(tmplDir, name))
}
