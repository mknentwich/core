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
	"regexp"
	"strings"
	"text/template"
)

const (
	tmplDir      = "templates"
	categoryTmpl = "category.*.tmpl"
	scoreTmpl    = "score.*.tmpl"
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
	generateCategories(rest.QueryCategoriesWithChildrenAndScores().([]database.Category), tmpl, scoreTmpl, outDir, ending, scoreEnding)
	return nil
}

func clean(dir string) error {
	files, err := ioutil.ReadDir(tmplDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := os.RemoveAll(path.Join(dir, file.Name())); err != nil {
			return err
		}
	}
}

func generateCategories(categories []database.Category, tmpl *template.Template, scoreTmpl *template.Template, outDir string, ending, scoreEnding string) {
	for _, category := range categories {
		generateCategory(&category, tmpl, scoreTmpl, outDir, ending, scoreEnding)
	}
}

func generateCategory(category *database.Category, tmpl *template.Template, scoreTmpl *template.Template, parentOutDir string, ending, scoreEnding string) {
	outDir := path.Join(parentOutDir, utils.SanitizePath(category.Name))
	generateCategories(category.Children, tmpl, scoreTmpl, outDir, ending, scoreEnding)
	for _, score := range category.Scores {
		generateScore(&score, scoreTmpl, outDir, scoreEnding)
	}
	file, err := os.OpenFile("_index."+ending, os.O_CREATE|os.O_RDWR, 0666)
	defer file.Close()
	if err == nil {
		err = tmpl.Execute(file, category)
	}
	if err != nil {
		log(context.LOG_ERROR, "cannot generate score %s.%s: %s", "_index", ending, err.Error())
	}
}

func generateScore(score *database.Score, tmpl *template.Template, parentOutDir string, ending string) {
	outFile := path.Join(parentOutDir, utils.SanitizePath(score.Title))
	file, err := os.OpenFile(outFile+"."+ending, os.O_CREATE|os.O_RDWR, 0666)
	defer file.Close()
	if err == nil {
		err = tmpl.Execute(file, score)
	}
	if err != nil {
		log(context.LOG_ERROR, "cannot generate score %s.%s: %s", outFile, ending, err.Error())
	}
}

func findTemplate(templateName string) (tmpl *template.Template, ending string, err error) {
	files, err := ioutil.ReadDir(tmplDir)
	if err != nil {
		return
	}
	regex, err := regexp.Compile(templateName)
	if err != nil {
		return
	}
	fileName := ""
	for _, file := range files {
		if regex.MatchString(file.Name()) && !file.IsDir() {
			if fileName != "" {
				err = &MultiMatchError{fmt.Sprintf("found multiple templates: %s and %s", fileName, file.Name())}
				return
			}
			fileName = file.Name()
		}
	}
	ending = templateName[strings.Index(templateName, "."):strings.LastIndex(templateName, ".")]
	tmpl, err = template.ParseFiles(fileName)
	return
}

type MultiMatchError struct {
	msg string
}

func (m *MultiMatchError) Error() string {
	return m.msg
}
