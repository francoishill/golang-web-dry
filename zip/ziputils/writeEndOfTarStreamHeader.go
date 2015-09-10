package ziputils

import (
	"archive/tar"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
)

const END_OF_TAR_FILENAME = "END_OF_TAR"

func writeEndOfTarStreamHeader(tarWriter *tar.Writer) {
	hdr := &tar.Header{
		Name: END_OF_TAR_FILENAME,
	}
	err := tarWriter.WriteHeader(hdr)
	CheckError(err)
}
