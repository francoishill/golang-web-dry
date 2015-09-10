package ziputils

import (
	"archive/tar"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/osutils"
	"net/http"
	"os"
)

func UploadFileToHttpResponseWriter(logger SimpleLogger, writer http.ResponseWriter, filePath string) {
	if !osutils.FileExists(filePath) {
		panic("File does not exist: " + filePath)
	}

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	file, err := os.OpenFile(filePath, 0, 0600)
	CheckError(err)
	defer file.Close()

	info, err := file.Stat()
	CheckError(err)

	writeFileToTarWriter(tarWriter, info, filePath, "", true)

	writeEndOfTarStreamHeader(tarWriter)
}
