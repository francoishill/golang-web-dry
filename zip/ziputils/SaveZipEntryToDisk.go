package ziputils

import (
	"archive/zip"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"os"
	"path/filepath"
)

func SaveZipEntryToDisk(logger SimpleLogger, destinationFolder string, fileEntry *zip.File) {
	rc, err := fileEntry.Open()
	CheckError(err)
	defer rc.Close()

	path := filepath.Join(destinationFolder, fileEntry.Name)
	if fileEntry.FileInfo().IsDir() {
		os.MkdirAll(path, fileEntry.Mode())
	} else {
		os.MkdirAll(filepath.Dir(path), fileEntry.Mode())

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileEntry.Mode())
		CheckError(err)
		defer file.Close()

		_, err = io.Copy(file, rc)
		CheckError(err)
	}
}
