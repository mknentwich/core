package media

import (
	"crypto/md5"
	"fmt"
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

//creates a md5sum from the filesystem
func readSumFromDisk(scoreId int, mediaType string) (string, error) {
	resPath := path.Join(outDir, mediaType, strconv.Itoa(scoreId))
	file, err := os.Open(resPath)
	if err != nil {
		return "", err
	}
	md := md5.New()
	_, err = io.Copy(md, file)
	if err != nil {
		return "", err
	}
	checksum := md.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

//removes a media from the filesystem
func removeMedia(scoreId int, mediaType string) error {
	resPath := path.Join(outDir, mediaType, strconv.Itoa(scoreId))
	return os.Remove(resPath)
}
