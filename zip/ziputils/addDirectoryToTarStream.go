package ziputils

import (
	"archive/tar"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"os"
	"path/filepath"
)

func addDirectoryToTarStream(w *tar.Writer, dir string, walkContext *dirWalkContext) {
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

		hdr, err := tar.FileInfoHeader(info, "")
		CheckError(err)
		hdr.Name = relPath

		err = w.WriteHeader(hdr)
		CheckError(err)

		if !info.IsDir() {
			file, err := os.Open(path)
			CheckError(err)
			defer file.Close()

			io.Copy(w, file)
		}

		return nil
	})

	CheckError(e)
}
