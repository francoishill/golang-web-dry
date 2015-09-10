package ziputils

import (
	"archive/zip"
	"bufio"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"os"
	"path/filepath"
)

func addDirectoryToZipStream(w *zip.Writer, dir string, walkContext *dirWalkContext) {
	e := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !walkContext.isMatch(info) {
			return nil
		}

		relPath := path[len(dir):]
		if relPath == "" {
			return nil
		}

		relPath = relPath[1:]
		file, err := os.Open(path)
		CheckError(err)
		defer file.Close()

		bufin := bufio.NewReader(file)

		zipEntryWriter, err := w.Create(relPath)
		CheckError(err)

		_, err = bufin.WriteTo(zipEntryWriter)
		CheckError(err)

		return nil
	})

	CheckError(e)
}
