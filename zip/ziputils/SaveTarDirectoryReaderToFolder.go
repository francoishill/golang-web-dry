package ziputils

import (
	"archive/tar"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"os"
	"path/filepath"
)

func SaveTarDirectoryReaderToFolder(logger SimpleLogger, bodyReader io.Reader, saveFolderPath string) {
	tarReader := tar.NewReader(bodyReader)

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		CheckError(err) //Check after checking for EOF

		relativePath := hdr.Name

		if hdr.FileInfo().IsDir() {
			fullDestinationDirPath := filepath.Join(saveFolderPath, relativePath)
			logger.Debug("(TAR) Creating directory %s", fullDestinationDirPath)
			os.MkdirAll(fullDestinationDirPath, os.FileMode(hdr.Mode))
			defer os.Chtimes(fullDestinationDirPath, hdr.AccessTime, hdr.ModTime)
		} else {
			fullDestinationFilePath := filepath.Join(saveFolderPath, relativePath)
			os.MkdirAll(filepath.Dir(fullDestinationFilePath), os.FileMode(hdr.Mode))

			logger.Debug("(TAR) Saving file %s", fullDestinationFilePath)
			file, err := os.OpenFile(fullDestinationFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
			CheckError(err)
			defer file.Close()
			defer os.Chtimes(fullDestinationFilePath, hdr.AccessTime, hdr.ModTime)

			_, err = io.Copy(file, tarReader)
			CheckError(err)

		}

		CheckError(err)
	}
}
