package ziputils

import (
	"archive/zip"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveZipDirectoryReaderToFolder(bodyReader io.Reader, saveFolderPath string) {
	tempDir := filepath.Join(os.TempDir(), "ZipDirs")
	err := os.MkdirAll(tempDir, 0600)
	CheckError(err)

	tempFile, err := ioutil.TempFile(tempDir, "networkzip-")
	CheckError(err)
	defer tempFile.Close()

	_, err = io.Copy(tempFile, bodyReader)
	CheckError(err)

	zipFile, err := zip.OpenReader(tempFile.Name())
	CheckError(err)
	defer zipFile.Close()

	for _, fileEntry := range zipFile.File {
		SaveZipEntryToDisk(saveFolderPath, fileEntry)
	}

	zipFile.Close()
	tempFile.Close()
	err = os.Remove(tempFile.Name())
	CheckError(err)
}
