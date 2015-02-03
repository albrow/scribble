package main

import (
	"github.com/OneOfOne/xxhash/native"
	"github.com/albrow/scribble/compilers"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/log"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
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
	log.Default.Println("Watching for changes...")
	if watcher == nil {
		var err error
		watcher, err = createWatcher()
		if err != nil {
			panic(err)
		}
	}

	// walk through source dir and watch all subdirectories
	// we have to do this because fsnotify is currently not recursive
	if err := watcher.Watch(config.SourceDir); err != nil {
		panic(err)
	}
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
			if err := watcher.Watch(path); err != nil {
				panic(err)
			}
			watchMutex.Unlock()
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

// createWatcher creates and returns an fsnotify.Watcher.
func createWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Process events
	go func() {
		defer util.Recovery(*compileTrace || *serveTrace)
		for {
			select {
			case ev := <-watcher.Event:
				if changed, err := fileDidChange(ev.Name); err != nil {
					panic(err)
				} else if changed {
					if err := compilers.FileChanged(ev.Name, *ev); err != nil {
						panic(err)
					}
				}
			case err := <-watcher.Error:
				panic(err)
			}
		}
	}()
	return watcher, nil
}

// fileDidChange uses the last known hash to determine whether or
// not the file actually changed. It solves the problem of false positives
// coming from fsnotify when used with a text editor that uses atomic saves.
func fileDidChange(path string) (bool, error) {
	if hash, found := fileHashes[path]; !found {
		// we have not hashed the file before.
		// hash it now and store the value
		if newHash, exists, err := calculateHashForPath(path); err != nil {
			return false, err
		} else if exists {
			fileHashes[path] = newHash
		}
		return true, nil
	} else {
		if newHash, exists, err := calculateHashForPath(path); err != nil {
			return false, err
		} else if !exists {
			// if the file no longer exists, it has been deleted
			// we should consider that a change and recompile
			delete(fileHashes, path)
			return true, nil
		} else if string(newHash) != string(hash) {
			// if the file does exist and has a different hash, there
			// was an actual change and we should recompile
			fileHashes[path] = newHash
			return true, nil
		}
		return false, nil
	}
}

// calculateHashForPath calculates a hash for the file at the given path.
// If the file does not exist, the second return value will be false.
func calculateHashForPath(path string) ([]byte, bool, error) {
	h := xxhash.New64()
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}
	if _, err := io.Copy(h, f); err != nil {
		return nil, false, err
	}
	if h.Size() == 0 {
		// The file existed, but it was empty
		return nil, true, nil
	}
	result := h.Sum(nil)
	return result, true, nil
}
