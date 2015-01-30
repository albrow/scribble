package main

import (
	"fmt"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/generators"
	"github.com/albrow/scribble/util"
	"os"
	"path/filepath"
)

// compile compiles all the contents of config.SourceDir and puts the compiled
// result in config.DestDir.
func compile(watch bool) {
	config.Parse()
	fmt.Println("--> compiling")
	if err := createDestDir(); err != nil {
		panic(err)
	}
	if err := removeAllOld(); err != nil {
		panic(err)
	}
	if err := generators.CompileAll(); err != nil {
		util.ChimeError(err)
	}
	if watch {
		watchAll()
	}
}

func createDestDir() error {
	if err := os.MkdirAll(config.DestDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			// If the directory already existed, that's fine,
			// otherwise return the error
			return err
		}
	}
	return nil
}

// removeAllOld removes all the files from config.DestDir
func removeAllOld() error {
	fmt.Println("    removing old files")
	// walk through the dest dir
	if err := filepath.Walk(config.DestDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == config.DestDir {
			// ignore the destDir itself
			return nil
		} else if info.IsDir() {
			// remove the dir and everything in it
			if err := os.RemoveAll(path); err != nil {
				if !os.IsNotExist(err) {
					// if the dir was already removed, that's fine.
					// if there was some other error, return it
					return err
				}
			}
			return filepath.SkipDir
		} else {
			// remove the file
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
