package osutils

import (
	"os"
)

func FileExists(file string) bool {
	stats, err := os.Stat(file)
	if err == nil {
		if !stats.IsDir() {
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
