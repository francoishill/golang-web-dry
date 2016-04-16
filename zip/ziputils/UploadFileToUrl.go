package ziputils

import (
	"archive/tar"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/francoishill/golang-web-dry/osutils"
)

func UploadFileToUrl(logger SimpleLogger, url, bodyType, filePath string, checkResponse func(resp *http.Response) error) {
	if !osutils.FileExists(filePath) {
		panic("File does not exist: " + filePath)
	}

	pipeReader, pipeWriter := io.Pipe()
	tarWriter := tar.NewWriter(pipeWriter)

	file, err := os.OpenFile(filePath, 0, 0600)
	CheckError(err)
	defer file.Close()

	info, err := file.Stat()
	CheckError(err)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	var goroutineErr error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				goroutineErr = fmt.Errorf("Cannot add directory to tar stream, error: %+v", r)
			}
		}()

		writeFileToTarWriter(tarWriter, info, filePath, "", true)
		writeEndOfTarStreamHeader(tarWriter)

		tarWriter.Close()
		pipeWriter.Close()

		wg.Done()
	}()

	resp, err := http.Post(url, bodyType, pipeReader)
	CheckError(err)
	resp.Body.Close()

	if checkResponse != nil {
		err := checkResponse(resp)
		CheckError(err)
	}

	wg.Wait()
	CheckError(goroutineErr)

	/*tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	file, err := os.OpenFile(filePath, 0, 0600)
	CheckError(err)
	defer file.Close()

	info, err := file.Stat()
	CheckError(err)

	writeFileToTarWriter(tarWriter, info, filePath, "", true)

	writeEndOfTarStreamHeader(tarWriter)*/
}
