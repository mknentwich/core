package media

import (
	"github.com/jinzhu/gorm"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"net/http"
)

func scoreExist(scoreId int) bool {
	return database.Receive().Find(&database.Score{
		Model: gorm.Model{ID: uint(scoreId)}}).Error == nil
}

func writeInternalError(err error, rw http.ResponseWriter) {
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log(context.LOG_ERROR, "%s", err.Error())
	}
}
