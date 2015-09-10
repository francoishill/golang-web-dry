package ziputils

import (
	"archive/tar"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/osutils"
	"io"
	"net/http"
	"sync"
)

func UploadDirectoryToUrl(logger SimpleLogger, url, bodyType, directoryPath string, walkContext *dirWalkContext, checkResponse func(resp *http.Response)) {
	if !osutils.DirectoryExists(directoryPath) {
		panic("Directory does not exist: " + directoryPath)
	}

	pipeReader, pipeWriter := io.Pipe()
	tarWriter := tar.NewWriter(pipeWriter)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		addDirectoryToTarStream(tarWriter, directoryPath, walkContext, true)
		tarWriter.Close()
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
