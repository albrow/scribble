package generators

import (
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
)

var Compilers = []Compiler{SassCompiler, AceCompiler}

type PathMatcher interface {
	// PathMatch returns a match string which will be
	// passed to filepath.Match to determine whether a
	// particular path should be ignored or added to the
	// PathMatcher
	PathMatch() string
	// IgnoreHidden returns true iff this PathMatcher
	// wants to ignore hidden system files, i.e. files
	// and directories which start with a '.'
	IgnoreHidden() bool
	// IgnoreUnderscore returns true iff this
	// PathMatcher wants to ignore files and directories
	// that start with an underscore.
	IgnoreUnderscore() bool
}

type Compiler interface {
	PathMatcher
	// Compile compiles a source file to an destination directory.
	// srcPath will be some path that matches according to the
	// PathMatcher. destDir will be the root destination directory.
	Compile(srcPath string, destDir string) error
	// CompileAll compiles all the files found in each path.
	// srcPaths will be all paths that match according to
	// the PathMatcher. destDir will be the root destination directory.
	CompileAll(srcPaths []string, destDir string) error
}

type Watcher interface {
	PathMatcher
	// PathChanged is triggered whenever a relevant file is changed
	// Typically, the Watcher should recompile certain files and put
	// them in the appropriate place in dest. srcPath will be some
	// path that matches according to the PathMatcher. destDir will be
	// the root destination directory.
	PathChanged(srcPath string, ev fsnotify.FileEvent, destDir string) error
}

// FindPaths iterates recursively through some root directory and
// returns all the matched paths for pm.
func FindPaths(root string, pm PathMatcher) ([]string, error) {
	paths := []string{}
	walkerFunc := createWalkerFunc(pm, &paths)
	if err := filepath.Walk(root, walkerFunc); err != nil {
		return nil, err
	}
	return paths, nil
}

// createWalkerFunc creates and returns a filepath.WalkerFunc which
// will check for matches, ignoring certain files and directories
// according to the PathMatcher, and append all the paths that match
// to paths.
func createWalkerFunc(pm PathMatcher, paths *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		firstChar := filepath.Base(path)[0]
		if pm.IgnoreHidden() && firstChar == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		} else if pm.IgnoreUnderscore() && firstChar == '_' {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		if matched, err := filepath.Match(pm.PathMatch(), filepath.Base(path)); err != nil {
			return err
		} else if matched {
			(*paths) = append(*paths, path)
		} else {
		}
		return nil
	}
}
