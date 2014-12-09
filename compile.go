package main

import (
	"fmt"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/generators"
	"github.com/albrow/scribble/util"
	"os"
	"path/filepath"
)

func compile(watch bool) {
	config.Parse()
	fmt.Println("--> compiling")
	createDestDir()
	removeAllOld()
	if err := generators.CompileAll(); err != nil {
		util.ChimeError(err)
	}
	// if watch {
	// 	watchAll()
	// }
}

func createDestDir() {
	if err := os.MkdirAll(config.DestDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			// If the directory already existed, that's fine,
			// otherwise panic.
			panic(err)
		}
	}
}

func removeAllOld() {
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
					// if the dir was already removed, that's fine
					// if there was some other error, panic
					panic(err)
				}
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
