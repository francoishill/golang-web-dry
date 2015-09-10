package ziputils

import (
	"archive/tar"
	"fmt"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io"
	"os"
)

func writeFileToTarWriter(tarWriter *tar.Writer, info os.FileInfo, absoluteFilePath, overwriteFileName string, isOnlyFile bool) {
	hdr, err := tar.FileInfoHeader(info, "")
	CheckError(err)
	if overwriteFileName != "" {
		hdr.Name = overwriteFileName
	}

	if hdr.Xattrs == nil {
		hdr.Xattrs = map[string]string{}
	}
	hdr.Xattrs["SIZE"] = fmt.Sprintf("%d", info.Size())
	if isOnlyFile {
		hdr.Xattrs["SINGLE_FILE_ONLY"] = "1"
	}

	err = tarWriter.WriteHeader(hdr)
	CheckError(err)

	if !info.IsDir() {
		file, err := os.Open(absoluteFilePath)
		CheckError(err)
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		CheckError(err)
	}
}
