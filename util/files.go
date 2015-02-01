package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateFileWithPath creates a file by first creating the directory
// the file will be placed in with os.MkdirAll (analogous to mkdir -p),
// and then creating the file itself. If the file already exists, it will
// overwrite the existing file. If there were any other problems creating
// the file, it will return an error.
func CreateFileWithPath(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	// First create the directory if needed
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			// If the dir already existed, that's fine.
			// For any other error, we should return it.
			return nil, err
		}
	}
	// Then create the file itself
	file, err := os.Create(path)
	if err != nil && !os.IsExist(err) {
		// If the file already existed, that's fine.
		// For any other error, we should return it.
		return nil, err
	}
	return file, nil
}

// CopyFile copies the file at srcePath to destPath. It creates any
// directories needed for destPath.
func CopyFile(srcPath string, destPath string) error {
	destFile, err := CreateFileWithPath(destPath)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}
	return nil
}

// RecursiveCopy copies everything from srcDir to destDir recursively.
// It is analogous to cp -R in unix systems.
func RecursiveCopy(srcDir string, destDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			destPath := strings.Replace(path, srcDir, destDir, 1)
			if err := CopyFile(path, destPath); err != nil {
				return err
			}
		}
		return nil
	})
}

// CreateEmptyFiles creates new, empty files for every path in paths. It
// does not write to them, and any old content that may have been there is
// erased.
func CreateEmptyFiles(paths []string) error {
	for _, path := range paths {
		if f, err := CreateFileWithPath(path); err != nil {
			return err
		} else {
			f.Close()
		}
	}
	return nil
}

// RemoveAllIfExists removes the directory identified by path if it exists.
// If it does not exist, calling this function has no effect. Contrary to
// the default behavior in the os package, RemoveAllIfExists will not return
// an error if path does not exist.
func RemoveAllIfExists(path string) error {
	if err := os.RemoveAll(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// RemoveIfExists removes the file identified by path if it exists. If it
// does not exist, calling this function has no effect. Contrary to the
// default behavior in the os package, RemoveIfExists will not return an
// error if path does not exist.
func RemoveIfExists(path string) error {
	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
