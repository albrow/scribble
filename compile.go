package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func compile(watch bool) {
	parseConfig()
	fmt.Println("--> compiling")
	fmt.Printf("    watch: %v\n", watch)
	removeOld()
	parsePosts()
	compileSass(watch)
	compilePages()
	compilePosts()
}

func removeOld() {
	fmt.Println("    removing old files")
	// walk through the dest dir
	if err := filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == destDir {
			// ignore the destDir itself
			return nil
		} else if info.IsDir() {
			if path == sassDestDir {
				// let sass handle this one
				return filepath.SkipDir
			}
			// remove the dir and everything in it
			if err := os.RemoveAll(path); err != nil {
				panic(err)
			}
			return filepath.SkipDir
		} else {
			// remove the file
			if err := os.Remove(path); err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
