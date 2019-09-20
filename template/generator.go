package template

import (
	"fmt"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/rest"
	"github.com/mknentwich/core/utils"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

const (
	tmplDir      = "templates"
	categoryTmpl = `category.*.tmpl`
	scoreTmpl    = `score.*.tmpl`
)

func Generate() error {
	tmpl, ending, err := findTemplate(categoryTmpl)
	if err != nil {
		return err
	}
	scoreTmpl, scoreEnding, err := findTemplate(scoreTmpl)
	if err != nil {
		return err
	}
	clean(outDir)
	generateCategories(rest.QueryCategoriesWithChildrenAndScores().([]database.Category), tmpl, scoreTmpl, outDir, ending, scoreEnding)
	return nil
}

func clean(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := os.RemoveAll(path.Join(dir, file.Name())); err != nil {
			return err
		}
	}
	return nil
}

func generateCategories(categories []database.Category, tmpl *template.Template, scoreTmpl *template.Template, outDir string, ending, scoreEnding string) {
	for _, category := range categories {
		generateCategory(&category, tmpl, scoreTmpl, outDir, ending, scoreEnding)
	}
}

func generateCategory(category *database.Category, tmpl *template.Template, scoreTmpl *template.Template, parentOutDir string, ending, scoreEnding string) {
	outDir := path.Join(parentOutDir, utils.SanitizePath(category.Name))
	err := os.Mkdir(outDir, 0700)
	generateCategories(category.Children, tmpl, scoreTmpl, outDir, ending, scoreEnding)
	for _, score := range category.Scores {
		generateScore(&score, scoreTmpl, outDir, scoreEnding)
	}
	file, err := os.OpenFile(path.Join(outDir, "_index")+ending, os.O_CREATE|os.O_RDWR, 0700)
	defer file.Close()
	if err == nil {
		err = tmpl.Execute(file, category)
	}
	if err != nil {
		log(context.LOG_ERROR, "cannot generate score %s%s: %s", "_index", ending, err.Error())
	}
}

func generateScore(score *database.Score, tmpl *template.Template, parentOutDir string, ending string) {
	outFile := path.Join(parentOutDir, utils.SanitizePath(score.Title))
	file, err := os.OpenFile(outFile+ending, os.O_CREATE|os.O_RDWR, 0700)
	defer file.Close()
	if err == nil {
		err = tmpl.Execute(file, score)
	}
	if err != nil {
		log(context.LOG_ERROR, "cannot generate score %s%s: %s", outFile, ending, err.Error())
	}
}

func findTemplate(templateName string) (tmpl *template.Template, ending string, err error) {
	files, err := ioutil.ReadDir(tmplDir)
	if err != nil {
		return
	}
	fileName := ""
	parts := strings.Split(templateName, "*")
	for _, file := range files {
		if len(file.Name()) >= len(templateName) && parts[0] == file.Name()[:len(parts[0])] && parts[1] == file.Name()[len(file.Name())-len(parts[1]):] && !file.IsDir() {
			if fileName != "" {
				err = &MultiMatchError{fmt.Sprintf("found multiple templates: %s and %s", fileName, file.Name())}
				return
			}
			fileName = file.Name()
		}
	}
	ending = fileName[strings.Index(fileName, "."):strings.LastIndex(fileName, ".")]
	tmpl, err = template.ParseFiles(path.Join(tmplDir, fileName))
	return
}

type MultiMatchError struct {
	msg string
}

func (m *MultiMatchError) Error() string {
	return m.msg
}
