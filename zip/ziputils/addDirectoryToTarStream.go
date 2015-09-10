package ziputils

import (
	"archive/tar"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"os"
	"path/filepath"
)

func addDirectoryToTarStream(tarWriter *tar.Writer, dir string, walkContext *dirWalkContext, writeEndHeader bool) {
	e := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !walkContext.isMatch(info) {
			return nil
		}

		relPath := path[len(dir):]
		if relPath == "" {
			return nil
		}

		relPath = relPath[1:]

		writeFileToTarWriter(tarWriter, info, path, relPath, false)

		return nil
	})

	CheckError(e)

	if writeEndHeader {
		writeEndOfTarStreamHeader(tarWriter)
	}
}
