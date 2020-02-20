package rest

import (
	"encoding/json"
	"github.com/mknentwich/core/context"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func Generate() error {
	err := clean(outDir)
	if err != nil {
		log(context.LOG_WARNING, "failed to wipe rest mirror directory: %s", err.Error())
	}
	err = mirrorQuery(QueryCategoriesFlat, context.Conf.RestMirror.CategoriesFlatPath)
	var sndErr error
	if sndErr = mirrorQuery(QueryCategoriesWithChildrenAndScores, context.Conf.RestMirror.CategoriesPath); err == nil {
		err = sndErr
	}
	if sndErr = mirrorQuery(QueryScoresFlat, context.Conf.RestMirror.ScoresPath); err == nil {
		err = sndErr
	}
	return nil
}

func mirrorQuery(query DataQuery, filename string) error {
	file, err := os.OpenFile(path.Join(outDir, filename), os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	return encoder.Encode(query())
}

//creates a worker for the template generator for every n minutes specified in the config
func worker() {
	ticker := time.NewTicker(time.Duration(context.Conf.RestMirror.Interval) * time.Minute)
	handleScheduleError(Generate())
	go func() {
		for range ticker.C {
			handleScheduleError(Generate())
		}
	}()
}

//TODO refactor with template
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

//handles a schedules error
//TODO refactor with template
func handleScheduleError(err error) {
	if err != nil {
		log(context.LOG_ERROR, "failed to mirrored rest api as scheduled: %s", err.Error())
	} else {
		log(context.LOG_INFO, "successfully mirrored rest api as scheduled")
	}
}
