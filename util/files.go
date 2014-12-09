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
// the file, it will panic, expecting that some caller higher up in the
// stack will reecover.
func CreateFileWithPath(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	// First create the directory if needed
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			// If the dir already existed, that's fine.
			// For any other error, we should panic.
			return nil, err
		}
	}
	// Then create the file itself
	file, err := os.Create(path)
	if err != nil && !os.IsExist(err) {
		// If the file already existed, that's fine.
		// For any other error, we should panic.
		return nil, err
	}
	return file, nil
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
