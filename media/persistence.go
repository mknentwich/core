package media

import (
	"io"
	"os"
	"path"
	"strconv"
)

//creates all necessary direcories
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

//writes a media to the filesystem
func saveMediaToDisk(scoreId int, mediaType string, reader io.Reader) error {
	resPath := path.Join(outDir, mediaType, strconv.Itoa(scoreId))
	file, err := os.OpenFile(resPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	return err
}

//reads a media from the filesystem
func readMediaFromDiskTo(scoreId int, mediaType string, writer io.Writer) error {
	resPath := path.Join(outDir, mediaType, strconv.Itoa(scoreId))
	file, err := os.Open(resPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

//removes a media from the filesystem
func removeMedia(scoreId int, mediaType string) error {
	resPath := path.Join(outDir, mediaType, strconv.Itoa(scoreId))
	return os.Remove(resPath)
}
