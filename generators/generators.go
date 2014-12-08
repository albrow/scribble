package generators

import (
	"github.com/albrow/scribble/config"
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
)

var Compilers = []Compiler{&PostsCompiler, &SassCompiler, &AceCompiler}

type Initer interface {
	// Init allows a Compiler or Watcher to do any necessary
	// setup before other methods are called. (e.g. set the
	// result of PathMatch based on some config variable).
	Init()
}

type Walker interface {
	// GetWalker returns a filepath.WalkFunc. The function returned
	// is responsible for adding any paths which the compiler/watcher
	// cares about to paths. The WalkFunc will be executed starting
	// at the source directory (lib.SrcDir) which is user-configurable.
	GetWalkFunc(paths *[]string) filepath.WalkFunc
}

type Compiler interface {
	Walker
	// Compile compiles a source file to an destination directory.
	// srcPath will be some path that matches according to the
	// Walker. destDir will be the root destination directory.
	Compile(srcPath string, destDir string) error
	// CompileAll compiles all the files found in each path.
	// srcPaths will be all paths that match according to
	// the Walker. destDir will be the root destination directory.
	CompileAll(srcPaths []string, destDir string) error
}

type Watcher interface {
	Walker
	// PathChanged is triggered whenever a relevant file is changed
	// Typically, the Watcher should recompile certain files and put
	// them in the appropriate place in dest. srcPath will be some
	// path that matches according to the PathMatcher. destDir will be
	// the root destination directory.
	PathChanged(srcPath string, ev fsnotify.FileEvent, destDir string) error
}

// FindPaths iterates recursively through some root directory and
// returns all the matched paths for w.
func FindPaths(root string, w Walker) ([]string, error) {
	paths := []string{}
	walkerFunc := w.GetWalkFunc(&paths)
	if err := filepath.Walk(root, walkerFunc); err != nil {
		return nil, err
	}
	return paths, nil
}

// CompileAll compiles all files in SrcDir by delegating each path to
// it's corresponding Compiler. If a path in SrcDir does not match any Compiler,
// it will be copied to DestDir directly. Any files or directories that start
// with an undercore ("_") will be ignored.
func CompileAll() error {
	for _, c := range Compilers {
		if initer, ok := c.(Initer); ok {
			// If the Compiler has an Init function, run it
			initer.Init()
		}
		paths, err := FindPaths(config.SourceDir, c)
		if err != nil {
			return err
		}
		if err := c.CompileAll(paths, config.DestDir); err != nil {
			return err
		}
	}
	return nil
}

// filenameMatchWalkFunc creates and returns a filepath.WalkFunc which
// will check for exact filename matches with pattern, according to the filepath.Match
// syntax, and append all the paths that match to paths. You can add options
// to ignore hidden files and directories (which start with a '.') or files
// which should typically be ignored by scribble (which start with a '_').
func filenameMatchWalkFunc(paths *[]string, pattern string, ignoreHidden bool, ignoreUnderscore bool) filepath.WalkFunc {
	return matchWalkFunc(paths, pattern, ignoreHidden, ignoreUnderscore, func(pattern, path string) (bool, error) {
		return filepath.Match(pattern, filepath.Base(path))
	})
}

// pathMatchWalkFunc creates and returns a filepath.WalkFunc which
// will check for full path matches with pattern, according to the filepath.Match
// syntax, and append all the paths that match to paths. You can add options
// to ignore hidden files and directories (which start with a '.') or files
// which should typically be ignored by scribble (which start with a '_').
func pathMatchWalkFunc(paths *[]string, pattern string, ignoreHidden bool, ignoreUnderscore bool) filepath.WalkFunc {
	return matchWalkFunc(paths, pattern, ignoreHidden, ignoreUnderscore, filepath.Match)
}

// matchWalkFunc creates and returns a filepath.WalkFunc which
// will check if a file path matches using matchFunc (i.e. when matchFunc returns true),
// and append all the paths that match to paths. You can add options to ignore hidden
// files and directories (which start with a '.') or files which should typically be
// ignored by scribble (which start with a '_').
func matchWalkFunc(paths *[]string, pattern string, ignoreHidden bool, ignoreUnderscore bool, matchFunc func(pattern, path string) (bool, error)) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		firstChar := filepath.Base(path)[0]
		if ignoreHidden && firstChar == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		} else if ignoreUnderscore && firstChar == '_' {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		if matched, err := matchFunc(pattern, path); err != nil {
			return err
		} else if matched {
			// fmt.Printf("%s matches %s\n", path, pattern)
			(*paths) = append(*paths, path)
		} else {
			// fmt.Printf("%s does not match %s\n", path, pattern)
		}
		return nil
	}
}
