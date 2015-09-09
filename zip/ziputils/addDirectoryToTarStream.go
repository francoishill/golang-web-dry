package ziputils

import (
	"archive/tar"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"os"
	"path/filepath"
)

func addDirectoryToTarStream(w *tar.Writer, dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
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

		hdr := &tar.Header{
			Name: relPath,
			Mode: 0600,
			Size: getFileSize(file),
		}
		err = w.WriteHeader(hdr)
		CheckError(err)

		io.Copy(w, file)

		return nil
	})
}
