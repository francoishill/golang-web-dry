package osutils

import (
	"os"
)

func DirectoryExists(dir string) bool {
	stats, err := os.Stat(dir)
	if err == nil {
		if stats.IsDir() {
			return true
		} else {
			return false
		}
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}
