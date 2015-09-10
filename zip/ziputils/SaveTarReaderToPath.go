package ziputils

import (
	"archive/tar"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"os"
	"path/filepath"
)

func SaveTarReaderToPath(logger SimpleLogger, bodyReader io.Reader, savePath string) {
	tarReader := tar.NewReader(bodyReader)

	foundEndOfTar := false
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		CheckError(err) //Check after checking for EOF

		if hdr.Name == END_OF_TAR_FILENAME {
			foundEndOfTar = true
			continue
		}

		relativePath := hdr.Name

		if hdr.FileInfo().IsDir() {
			fullDestinationDirPath := filepath.Join(savePath, relativePath)
			logger.Debug("(TAR) Creating directory %s", fullDestinationDirPath)
			os.MkdirAll(fullDestinationDirPath, os.FileMode(hdr.Mode))
			defer os.Chtimes(fullDestinationDirPath, hdr.AccessTime, hdr.ModTime)
		} else {
			fullDestinationFilePath := filepath.Join(savePath, relativePath)
			if val, ok := hdr.Xattrs["SINGLE_FILE_ONLY"]; ok && val == "1" {
				fullDestinationFilePath = savePath
			}

			os.MkdirAll(filepath.Dir(fullDestinationFilePath), os.FileMode(hdr.Mode))

			logger.Debug("(TAR) Saving file %s", fullDestinationFilePath)
			file, err := os.OpenFile(fullDestinationFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
			CheckError(err)
			defer file.Close()
			defer os.Chtimes(fullDestinationFilePath, hdr.AccessTime, hdr.ModTime)

			_, err = io.Copy(file, tarReader)
			CheckError(err)
		}
	}

	if !foundEndOfTar {
		panic("TAR stream validation failed, something has gone wrong during the transfer.")
	}
}
