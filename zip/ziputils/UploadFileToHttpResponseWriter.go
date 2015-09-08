package ziputils

import (
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"net/http"
	"os"
)

func UploadFileToHttpResponseWriter(writer http.ResponseWriter, filePath string) {
	file, err := os.OpenFile(filePath, 0, 0600)
	CheckError(err)

	writer.Header().Set("Content-Length", fmt.Sprintf("%d", getFileSize(file)))

	_, err = io.Copy(writer, file)
	CheckError(err)
}
