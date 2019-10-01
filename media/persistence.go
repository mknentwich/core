package media

import (
	"io"
	"os"
	"path"
	"strconv"
)

func createDirs() error {
	var err error
	for _, mediaType := range mediaTypes {
		err = os.MkdirAll(path.Join(outDir, mediaType), 0700)
		if err != nil {
			return err
		}
	}
	return err
}

func saveMediaToDisk(scoreId int, mediaType string, reader io.Reader) error {
	resPath := path.Join(outDir, mediaType, strconv.Itoa(scoreId))
	file, err := os.OpenFile(resPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	return err
}
