package ziputils

import (
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/osutils"
	"io"
	"net/http"
	"os"
)

func UploadFileToHttpResponseWriter(logger SimpleLogger, writer http.ResponseWriter, filePath string) {
	if !osutils.FileExists(filePath) {
		panic("File does not exist: " + filePath)
	}

	file, err := os.OpenFile(filePath, 0, 0600)
	CheckError(err)

	writer.Header().Set("Content-Length", fmt.Sprintf("%d", getFileSize(file)))

	_, err = io.Copy(writer, file)
	CheckError(err)
}
