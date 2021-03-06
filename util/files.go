// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

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
			if err := f.Close(); err != nil {
				return err
			}
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

// RemoveEmptyDirs recursively iterates through path and removes any empty
// directories within it.
func RemoveEmptyDirs(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Open the directory to see if it's empty
			file, err := os.Open(path)
			if err != nil {
				switch err {
				case os.ErrNotExist:
					// If the directory we were going to maybe delete doesn't exist
					// anymore, that's fine
					return nil
				default:
					// If there was some other error, return it
					return err
				}
			}

			if _, err := file.Readdirnames(1); err != nil {
				switch err {
				case io.EOF:
					// This means the directory has no files, we need to delete it
					if err := RemoveAllIfExists(path); err != nil {
						return err
					}
				default:
					// If there was some other error, return it
					return err
				}
			}
		}
		return nil
	})
}
