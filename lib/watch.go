package lib

import (
	"fmt"
	"github.com/OneOfOne/xxhash/native"
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

func watchAll() {
	fmt.Println("--> watching for changes")
	if watcher == nil {
		watcher = createWatcher()
	}
	// TODO: be more intelligent here. E.g. if a sass file changes,
	// only recompile the sass files.

	// walk through source dir and watch all subdirectories
	// we have to do this because fsnotify is currently not recursive
	watcher.Watch(SourceDir)
	if err := filepath.Walk(SourceDir, func(path string, info os.FileInfo, err error) error {
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
		defer Recovery()
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

// func reactToPostChange(ev *fsnotify.FileEvent) {
// 	if ev.IsCreate() || ev.IsModify() {
// 		// recompile the single post
// 		p := getOrCreatePostFromPath(ev.Name)
// 		p.parse()
// 		p.compile()
// 		// recompile all pages, since they may depend on posts
// 		compilePages()
// 	} else if ev.IsRename() {
// 		// recompile all posts
// 		// TODO: detect which post was renamed and only recompile that one?
// 		if p := getPostByPath(ev.Name); p != nil {
// 			p.remove()
// 		}
// 		parsePosts()
// 		compilePosts()
// 		compilePages()
// 	} else if ev.IsDelete() {
// 		if p := getPostByPath(ev.Name); p != nil {
// 			p.remove()
// 		}
// 		// recompile all pages, since they may depend on posts
// 		compilePages()
// 	}
// }

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

// calculates a hash for the file at the given path. If the file
// does not exist, the second return value will be false.
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
