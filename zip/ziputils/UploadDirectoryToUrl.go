package ziputils

import (
	"archive/zip"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"net/http"
	"sync"
)

func UploadDirectoryToUrl(url, bodyType, directoryPath string) {
	pipeReader, pipeWriter := io.Pipe()
	zipWriter := zip.NewWriter(pipeWriter)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		addDirectoryToZipStream(zipWriter, directoryPath)
		zipWriter.Close()
		pipeWriter.Close()

		wg.Done()
	}()

	resp, err := http.Post(url, bodyType, pipeReader)
	CheckError(err)
	resp.Body.Close()

	wg.Wait()
}
