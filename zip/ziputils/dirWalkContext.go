package ziputils

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"os"
	"path/filepath"
)

type dirWalkContext struct {
	FileFilterPattern string
}

func (d *dirWalkContext) isMatch(info os.FileInfo) bool {
	if d.FileFilterPattern == "" {
		//No filter
		return true
	}
	if info.IsDir() {
		//Always let
		return true
	}

	matched, err := filepath.Match(d.FileFilterPattern, info.Name())
	CheckError(err)
	return matched
}

func (d *dirWalkContext) DeleteDirectory(dir string) {
	if d.FileFilterPattern == "" {
		err := os.RemoveAll(dir)
		CheckError(err)
	} else {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() &&
				d.isMatch(info) {
				//Skip directories if we are filtering for files
				return os.Remove(path)
			}

			return nil
		})
	}
}

/*
Creates a new instance of dirWalkContext.

For example to only filter and find .txt files, use:
	wc := NewDirWalkContext("*.txt")
*/
func NewDirWalkContext(fileFilterPattern string) *dirWalkContext {
	return &dirWalkContext{
		FileFilterPattern: fileFilterPattern,
	}
}
