package ziputils

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"os"
	"path/filepath"

	"github.com/francoishill/golang-web-dry/osutils"
)

func SaveReaderToFile(logger SimpleLogger, bodyReader io.Reader, saveFilePath string) {
	fullDestinationDirPath := filepath.Dir(saveFilePath)
	if !osutils.DirectoryExists(fullDestinationDirPath) {
		logger.Debug("(TAR) Creating directory '%s' ( parent of file '%s')", fullDestinationDirPath, filepath.Base(saveFilePath))
		err := os.MkdirAll(fullDestinationDirPath, os.FileMode(0655))
		CheckError(err)
	}

	out, err := os.Create(saveFilePath)
	CheckError(err)
	defer out.Close()

	_, err = io.Copy(out, bodyReader)
	CheckError(err)
}
