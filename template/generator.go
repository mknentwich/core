package template

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/rest"
	"github.com/mknentwich/core/utils"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
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
	err = clean(outDir)
	if err != nil {
		log(context.LOG_WARNING, "failed to wipe template directory: %s", err.Error())
	}
	generateCategories(rest.QueryCategoriesWithChildrenAndScores().([]database.Category), tmpl, scoreTmpl, outDir, ending, scoreEnding)
	rootCategory := database.Category{Name: "Übersicht", Model: gorm.Model{ID: 0}}
	generateCategory(&rootCategory, tmpl, scoreTmpl, outDir, ending, scoreEnding)
	return nil
}

//creates a worker for the template generator for every n minutes specified in the config
func worker() {
	ticker := time.NewTicker(time.Duration(context.Conf.TemplateInterval) * time.Minute)
	handleScheduleError(Generate())
	go func() {
		for range ticker.C {
			handleScheduleError(Generate())
		}
	}()
}

//handles a schedules error
func handleScheduleError(err error) {
	if err != nil {
		log(context.LOG_ERROR, "failed to generated scheduled templates: %s", err.Error())
	} else {
		log(context.LOG_INFO, "successfully generated scheduled templates")
	}
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
	if category.ID == 0 {
		outDir = parentOutDir
	}
	err := os.Mkdir(outDir, 0700)
	generateCategories(category.Children, tmpl, scoreTmpl, outDir, ending, scoreEnding)
	for _, score := range category.Scores {
		(&score).Category = category
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
	tmpl, err = template.New(fileName).Funcs(template.FuncMap{
		"sanitize": utils.SanitizePath}).ParseFiles(path.Join(tmplDir, fileName))
	ending = fileName[strings.Index(fileName, "."):strings.LastIndex(fileName, ".")]
	if tmpl != nil {
		tmpl = tmpl.Funcs(template.FuncMap{
			"sanitize": utils.SanitizePath})
	}
	return
}

type MultiMatchError struct {
	msg string
}

func (m *MultiMatchError) Error() string {
	return m.msg
}
