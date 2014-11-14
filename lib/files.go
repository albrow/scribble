package lib

import (
	"os"
	"path/filepath"
)

// createFileWithPath creates a file by first creating the directory
// the file will be placed in with os.MkdirAll (analogous to mkdir -p),
// and then creating the file itself. If the file already exists, it will
// overwrite the existing file. If there were any other problems creating
// the file, it will panic, expecting that some caller higher up in the
// stack will reecover.
func createFileWithPath(path string) *os.File {
	dir := filepath.Dir(path)
	// First create the directory if needed
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			// If the dir already existed, that's fine.
			// For any other error, we should panic.
			panic(err)
		}
	}
	// Then create the file itself
	file, err := os.Create(path)
	if err != nil && !os.IsExist(err) {
		// If the file already existed, that's fine.
		// For any other error, we should panic.
		panic(err)
	}
	return file
}
