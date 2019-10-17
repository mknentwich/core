package media

import (
	"github.com/jinzhu/gorm"
	"github.com/mknentwich/core/database"
	"os"
	"path"
	"strconv"
)

func scoreExist(scoreId int, mediaType string) bool {
	return database.Receive().Find(&database.Score{
		Model: gorm.Model{ID: uint(scoreId)}}).Error == nil && func() bool {
		_, err := os.Stat(path.Join(outDir, mediaType, strconv.Itoa(scoreId)))
		return err == nil
	}()
}
