package generators

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
	"strings"
)

// Compilers is a slice of all known Compilers.
// NOTE: it is important that PostsCompiler is the first
// item in the slice, because some other compilers rely on
// the existence of a list of parsed Post objects. For example,
// AceCompiler relies on the Posts function returning the correct
// results inside of ace templates.
var Compilers = []Compiler{&PostsCompiler, &SassCompiler, &AceCompiler}

// CompilerPaths is a map of Compiler to the matched paths for that Compiler
var CompilerPaths = map[Compiler][]string{}

// UnmatchedPaths is a slice of paths which do not match any compiler
var UnmatchedPaths = []string{}

// MatchFunc represents a function which should return true iff
// path matches some pattern. Compilers and Watchers return a MatchFunc
// to specify which paths they are concerned with.
type MatchFunc func(path string) (bool, error)

// Initer is an interface satisfied by any Compiler which needs
// to do something before Compile or CompileAll are called.
type Initer interface {
	// Init allows a Compiler or Watcher to do any necessary
	// setup before other methods are called. (e.g. set the
	// result of PathMatch based on some config variable). The
	// Init method is not required, but it will be called if
	// it exists.
	Init()
}

// PathMatcher is responsible for managing a certain set of files
// selected via a simple matching function.
type PathMatcher interface {
	// GetMatchFunc returns a MatchFunc, which in this context
	// will be used by a Compiler or Watcher to specify which paths
	// it is concerned with.
	GetMatchFunc() MatchFunc
}

// Compiler is capable of compiling a certain type of file.
type Compiler interface {
	PathMatcher
	// Compile compiles a source file identified by srcPath.
	// srcPath will be some path that matches according to the
	// MatchFunc for the Compiler.
	Compile(srcPath string) error
	// CompileAll compiles all the files found in each path.
	// srcPaths will be all paths that match according to
	// the MatchFunc for the Compiler.
	CompileAll(srcPaths []string) error
}

// Watcher is responsible for watching a specific set of files and
// reacting to changes to those files.
type Watcher interface {
	PathMatcher
	// PathChanged is triggered whenever a relevant file is changed
	// Typically, the Watcher should recompile certain files.
	// srcPath will be some path that matches according to the MatchFunc
	// for the Watcher. ev is the FileEvent associated with the change.
	PathChanged(srcPath string, ev fsnotify.FileEvent) error
}

// FindPaths iterates recursively through config.SourceDir and
// returns all the matched paths for m.
func FindPaths(m PathMatcher) ([]string, error) {
	paths := []string{}
	matchFunc := m.GetMatchFunc()
	walkFunc := matchWalkFunc(&paths, matchFunc)
	if err := filepath.Walk(config.SourceDir, walkFunc); err != nil {
		return nil, err
	}
	return paths, nil
}

// CompileAll compiles all files in config.SourceDir by delegating each path to
// it's corresponding Compiler. If a path in config.SourceDir does not match any Compiler,
// it will be copied to config.DestDir directly.
func CompileAll() error {
	initCompilers()
	if err := delegatePaths(); err != nil {
		return err
	}
	for _, c := range Compilers {
		paths, found := CompilerPaths[c]
		if found && len(paths) > 0 {
			if err := c.CompileAll(paths); err != nil {
				return err
			}
		}
	}
	if err := copyUnmatchedPaths(UnmatchedPaths); err != nil {
		return err
	}
	return nil
}

// initCompilers calls the Init method for each compiler that has it
func initCompilers() {
	for _, c := range Compilers {
		if initer, ok := c.(Initer); ok {
			// If the Compiler has an Init function, run it
			initer.Init()
		}
		CompilerPaths[c] = []string{}
	}
}

// delegatePaths walks through the source directory, checks if a path matches according
// to the MatchFunc for each compiler, and adds the path to CompilerPaths if it does
// match.
func delegatePaths() error {
	return filepath.Walk(config.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		matched := false
		for _, c := range Compilers {
			if match, err := c.GetMatchFunc()(path); err != nil {
				return err
			} else if match {
				matched = true
				CompilerPaths[c] = append(CompilerPaths[c], path)
			}
		}
		if !matched && !info.IsDir() {
			// If the path didn't match any compilers according to their MatchFuncs,
			// add it to the list of unmatched paths. These will be copied from config.SourceDir
			// to config.DestDir without being changed.
			UnmatchedPaths = append(UnmatchedPaths, path)
		}
		return nil
	})
}

// copyUnmatchedPaths copies paths from config.SourceDir to config.DestDir without changing them. It perserves
// directory structures, so e.g., source/archive/index.html becomes public/archive/index.html.
func copyUnmatchedPaths(paths []string) error {
	for _, path := range paths {
		destPath := strings.Replace(path, config.SourceDir, config.DestDir, 1)
		if err := util.CopyFile(path, destPath); err != nil {
			return err
		}
	}
	return nil
}

// filenameMatchFunc creates and returns a MatchFunc which
// will check for exact matches between the filename for some path (i.e. filepath.Base(path))
// and pattern, according to the filepath.Match semantics. You can add options to ignore hidden
// files and directories (which start with a '.') or files which should typically be ignored by
// scribble (which start with a '_').
// BUG: Filename matching is not expected to work on windows.
func filenameMatchFunc(pattern string, ignoreHidden bool, ignoreUnderscore bool) MatchFunc {
	return func(path string) (bool, error) {
		if ignoreHidden {
			// Check for hidden files and directories, i.e. those
			// that begin with a '.'. If we find the substring "/."
			// it must mean that some file or directory in the path
			// is hidden.
			// TODO: Make this compatible with windows.
			if strings.Contains(path, "/.") {
				return false, nil
			}
		}
		if ignoreUnderscore {
			// Check for files and directories that begin with a '_',
			// which have special meaning in scribble and should typically
			// be ignored. If we find the substring "/_" it must mean that some
			// file or directory in the path starts with an underscore.
			// TODO: Make this compatible with windows.
			if strings.Contains(path, "/_") {
				return false, nil
			}
		}
		return filepath.Match(pattern, filepath.Base(path))
	}
}

// pathMatchFunc creates and returns a MatchFunc which
// will check for exact matches between some full path and pattern, according to the
// filepath.Match semantics. You can add options to ignore hidden files and directories
// (which start with a '.') or files which should typically be ignored by scribble (which
// start with a '_').
// BUG: Path matching is not expected to work on windows.
func pathMatchFunc(pattern string, ignoreHidden bool, ignoreUnderscore bool) MatchFunc {
	return func(path string) (bool, error) {
		if ignoreHidden {
			// Check for hidden files and directories, i.e. those
			// that begin with a '.'. If we find the substring "/."
			// it must mean that some file or directory in the path
			// is hidden.
			// TODO: Make this compatible with windows.
			if strings.Contains(path, "/.") {
				return false, nil
			}
		}
		if ignoreUnderscore {
			// Check for files and directories that begin with a '_',
			// which have special meaning in scribble and should typically
			// be ignored. If we find the substring "/_" it must mean that some
			// file or directory in the path starts with an underscore.
			// TODO: Make this compatible with windows.
			if strings.Contains(path, "/_") {
				return false, nil
			}
		}
		return filepath.Match(pattern, path)
	}
}

// matchWalkFunc creates and returns a filepath.WalkFunc which
// will check if a file path matches using matchFunc (i.e. when matchFunc returns true),
// and append all the paths that match to paths. Typically, this should only be used when
// you want to get the paths for a specific Compiler/Watcher and not for any of the others,
// e.g. for testing.
func matchWalkFunc(paths *[]string, matchFunc func(path string) (bool, error)) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if matched, err := matchFunc(path); err != nil {
			return err
		} else if matched {
			// fmt.Printf("%s matches %s\n", path, pattern)
			if paths != nil {
				(*paths) = append(*paths, path)
			}
		} else {
			// fmt.Printf("%s does not match %s\n", path, pattern)
		}
		return nil
	}
}
