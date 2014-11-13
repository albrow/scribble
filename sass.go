package main

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func compileSass() {
	fmt.Println("    compiling sass")

	// walk through the sass source dir
	if err := filepath.Walk(sassSourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == sassSourceDir {
			// don't do anything with the sassSourceDir itself
			return nil
		}
		base := info.Name()
		if base[0] == '_' {
			// ignore hidden files
			if info.IsDir() {
				// skip any files in directories that start with '_'
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".scss" {
			compileSassFromPath(path)
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func compileSassFromPath(path string) {
	// parse path and figure out destPath
	destDir := strings.Replace(filepath.Dir(path), sassSourceDir, sassDestDir, 1)
	srcFile := filepath.Base(path)
	destFile := strings.Replace(srcFile, ".scss", ".css", 1)
	destPath := fmt.Sprintf("%s/%s", destDir, destFile)
	color.Printf("@g    CREATE: %s -> %s\n", path, destPath)

	// create the destDir if needed
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		// if the dir already exists, that's fine
		// if there was some other error, panic
		if !os.IsExist(err) {
			panic(err)
		}
	}

	// set up and execute the command, capturing the output only if there was an error
	cmd := exec.Command("sassc", path, destPath)
	response, err := cmd.CombinedOutput()
	if err != nil {
		chimeError(string(response))
	}
}
