package ziputils

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"os"
)

func SaveReaderToFile(bodyReader io.Reader, saveFilePath string) {
	out, err := os.Create(saveFilePath)
	CheckError(err)
	defer out.Close()

	_, err = io.Copy(out, bodyReader)
	CheckError(err)
}
