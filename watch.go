package main

import (
	"fmt"
	"github.com/OneOfOne/xxhash/native"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"github.com/wsxiaoys/terminal/color"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var watchedPaths = []string{}
var watcher *fsnotify.Watcher
var watchMutex = sync.Mutex{}

// a map of known file hashes. This is required to
// determine whether a file actually changed. It is a workaround
// to fix the bug that occurs when a text editor uses atomic saves,
// which triggers multiple watch events even though the file was
// only saved once.
var fileHashes = map[string][]byte{}

// watchAll begins watching all the files in config.SourceDir and reacts
// to any changes.
func watchAll() {
	fmt.Println("--> watching for changes")
	if watcher == nil {
		watcher = createWatcher()
	}
	// TODO: be more intelligent here. E.g. if a sass file changes,
	// only recompile the sass files.

	// walk through source dir and watch all subdirectories
	// we have to do this because fsnotify is currently not recursive
	watcher.Watch(config.SourceDir)
	if err := filepath.Walk(config.SourceDir, func(path string, info os.FileInfo, err error) error {
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

// createWatcher creates and returns an fsnotify.Watcher.
func createWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	// Process events
	go func() {
		defer util.Recovery()
		for {
			select {
			case ev := <-watcher.Event:
				base := filepath.Base(ev.Name)
				if base[0] == '.' {
					// ignore hidden system files
					continue
				}
				if fileDidChange(ev.Name) {
					color.Printf("@y    CHANGED: %s\n", ev.Name)
					// TODO: rewrite this
				}
			case err := <-watcher.Error:
				panic(err)
			}
		}
	}()
	return watcher
}

// fileDidChange uses the last known hash to determine whether or
// not the file actually changed. It solves the problem of false positives
// coming from fsnotify when used with a text editor that uses atomic saves.
func fileDidChange(path string) bool {
	if hash, found := fileHashes[path]; !found {
		// we have not hashed the file before.
		// hash it now and store the value
		newHash, exists := calculateHashForPath(path)
		if exists {
			fileHashes[path] = newHash
		}
		return true
	} else {
		newHash, exists := calculateHashForPath(path)
		if !exists {
			// if the file no longer exists, it has been deleted
			// we should consider that a change and recompile
			delete(fileHashes, path)
			return true
		} else if string(newHash) != string(hash) {
			// if the file does exist and has a different hash, there
			// was an actual change and we should recompile
			fileHashes[path] = newHash
			return true
		}
		return false
	}
}

// calculateHashForPath calculates a hash for the file at the given path.
// If the file does not exist, the second return value will be false.
func calculateHashForPath(path string) ([]byte, bool) {
	h := xxhash.New64()
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false
		} else {
			panic(err)
		}
	}
	io.Copy(h, f)
	result := h.Sum(nil)
	return result, true
}
