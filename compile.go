package main

import (
	"github.com/albrow/scribble/compilers"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/log"
	"github.com/albrow/scribble/util"
	"os"
	"path/filepath"
)

// compile compiles all the contents of config.SourceDir and puts the compiled
// result in config.DestDir.
func compile(watch bool) {
	config.Parse()
	log.Default.Println("Compiling...")
	if err := createDestDir(); err != nil {
		panic(err)
	}
	if err := removeAllOld(); err != nil {
		panic(err)
	}
	if err := compilers.CompileAll(); err != nil {
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
	log.Default.Println("Removing old files...")
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
			if err := util.RemoveAllIfExists(path); err != nil {
				return err
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
