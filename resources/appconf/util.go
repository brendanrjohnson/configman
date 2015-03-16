package appconf

import (
	"os"
)

// isFileExist reports whether path exists.
func isFileExist(fpath string) bool {
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return false
	}
	return true
}
