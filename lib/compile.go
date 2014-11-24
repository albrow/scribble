package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

type Context map[string]interface{}

var context = Context{}

func Compile(watch bool) {
	parseConfig()
	fmt.Println("--> compiling")
	removeAllOld()
	// parsePosts()

	if watch {
		watchAll()
	}
}

func GetContext() Context {
	return context
}

func removeAllOld() {
	fmt.Println("    removing old files")
	// walk through the dest dir
	if err := filepath.Walk(DestDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == DestDir {
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
