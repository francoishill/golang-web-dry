package ziputils

import (
	"archive/zip"
	"net/http"
)

func UploadDirectoryToHttpResponseWriter(writer http.ResponseWriter, directoryPath string) {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	addDirectoryToZipStream(zipWriter, directoryPath)
}
