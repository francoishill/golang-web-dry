package ziputils

import (
	"archive/zip"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"net/http"
	"sync"
)

func UploadDirectoryToUrl(url, bodyType, directoryPath string, checkResponse func(resp *http.Response)) {
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

	if checkResponse != nil {
		checkResponse(resp)
	}

	wg.Wait()
}
