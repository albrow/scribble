package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"path/filepath"
	"sync"
)

var watchedPaths = []string{}
var watcher *fsnotify.Watcher
var watchMutex = sync.Mutex{}

func watchAll() {
	fmt.Println("--> watching for changes")
	if watcher == nil {
		watcher = createWatcher()
	}
	// TODO: be more intelligent here. E.g. if a sass file changes,
	// only recompile the sass files.

	// walk through source dir and watch all subdirectories
	// we have to do this because fsnotify is currently not recursive
	watcher.Watch(sourceDir)
	if err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name()[0] == '.' {
			// ignore hidden system files
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			watchMutex.Lock()
			watchedPaths = append(watchedPaths, path)
			watcher.Watch(path)
			watchMutex.Unlock()
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func createWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	// Process events
	go func() {
		defer recovery()
		for {
			select {
			case ev := <-watcher.Event:
				base := filepath.Base(ev.Name)
				if base[0] == '.' {
					// ignore hidden system files
					continue
				}
				fmt.Printf("ev: %+v\n", ev)
				color.Printf("@y    CHANGED: %s\n", ev.Name)
				compile(false)
			case err := <-watcher.Error:
				panic(err)
			}
		}
	}()
	return watcher
}
