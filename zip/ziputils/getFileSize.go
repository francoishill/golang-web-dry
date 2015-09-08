package ziputils

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"os"
)

func getFileSize(file *os.File) int64 {
	fi, err := file.Stat()
	CheckError(err)
	return fi.Size()
}
