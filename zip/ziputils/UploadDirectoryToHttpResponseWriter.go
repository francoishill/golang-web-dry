package ziputils

import (
	"archive/tar"
	"github.com/francoishill/golang-web-dry/osutils"
	"net/http"
)

func UploadDirectoryToHttpResponseWriter(logger SimpleLogger, writer http.ResponseWriter, directoryPath string, walkContext *dirWalkContext) {
	if !osutils.DirectoryExists(directoryPath) {
		panic("Directory does not exist: " + directoryPath)
	}

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	addDirectoryToTarStream(tarWriter, directoryPath, walkContext, true)
}
