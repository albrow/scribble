package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func compileSass(watch bool) {
	// TODO: implement watch
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
	fmt.Printf("    %s -> %s\n", path, destPath)

	// create the destDir if needed
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		// if the dir already exists, that's fine
		// if there was some other error, panic
		if !os.IsExist(err) {
			panic(err)
		}
	}

	// set up and execute the command, piping output to stdout
	cmd := exec.Command("sassc", path, destPath)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
