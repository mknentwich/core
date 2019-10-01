package media

import (
	"os"
	"path"
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
