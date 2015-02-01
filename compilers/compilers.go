package compilers

import (
	"fmt"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"github.com/wsxiaoys/terminal/color"
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
var Compilers = []Compiler{&PostsCompiler, &SassCompiler, &HtmlTemplatesCompiler}

// CompilerPaths is a map of Compiler to the matched paths for that Compiler
var CompilerPaths = map[Compiler][]string{}

// UnmatchedPaths is a slice of paths which do not match any compiler
var UnmatchedPaths = []string{}

// noHiddenNoIgnore is a MatchFunc which returns true for any path that is
// does not begin with a "." or "_" and is not inside any directory which begins
// with a "." or "_".
var noHiddenNoIgnore = filenameMatchFunc("*", true, true)

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

// Compiler is capable of compiling a certain type of file. It
// also is responsible for watching for changes to certain types
// of files.
type Compiler interface {
	// CompileMatchFunc returns a MatchFunc which will be applied
	// to every path in config.SourceDir to determine which paths
	// a Compiler is responsible for compiling.
	CompileMatchFunc() MatchFunc
	// Compile compiles a source file identified by srcPath.
	// srcPath will be some path that matches according to the
	// MatchFunc for the Compiler.
	Compile(srcPath string) error
	// CompileAll compiles all the files found in each path.
	// srcPaths will be all paths that match according to
	// the MatchFunc for the Compiler.
	CompileAll(srcPaths []string) error
	// RemoveAllOld removes all files which this compiler has created
	// in config.DestDir. A Compiler is responsible for keeping track
	// of the files it has created and removing them when this method
	// is called.
	RemoveOld() error
	// WatchMatchFunc returns a MatchFunc which will be applied
	// to every path in config.SourceDir to determine which paths
	// a Compiler is responsible for watching. Note that the files
	// that are watched may not be the same as those that are compiled.
	// E.g, files that start with an underscore are typically not compiled,
	// but may be imported or used by other files that are compiled, and
	// therefore should be watched.
	WatchMatchFunc() MatchFunc
	// FileChanged is triggered whenever a relevant file is changed.
	// Typically, the Compiler should recompile certain files.
	// srcPath will be some path that matches according to WatchMatchFunc,
	// and ev is the FileEvent associated with the change.
	FileChanged(srcPath string, ev fsnotify.FileEvent) error
}

// FindPaths iterates recursively through config.SourceDir and
// returns all the matched paths using mf as a MatchFunc.
func FindPaths(mf MatchFunc) ([]string, error) {
	paths := []string{}
	walkFunc := matchWalkFunc(&paths, mf)
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
	if err := delegateCompilePaths(); err != nil {
		return err
	}
	for _, c := range Compilers {
		if err := compileAllForCompiler(c); err != nil {
			return err
		}
	}
	if err := copyUnmatchedPaths(UnmatchedPaths); err != nil {
		return err
	}
	return nil
}

// compileAllForCompiler recompiles all paths that are matched according to the given compiler's
// MatchFunc
func compileAllForCompiler(c Compiler) error {
	paths, found := CompilerPaths[c]
	if found && len(paths) > 0 {
		if err := c.CompileAll(paths); err != nil {
			return err
		}
	}
	return nil
}

// FileChanged delegates file changes to the appropriate compiler. If srcPath does not match any
// Compiler, it will be copied to config.DestDir directly.
func FileChanged(srcPath string, ev fsnotify.FileEvent) error {
	hasMatch := false
	for _, c := range Compilers {
		if match, err := c.WatchMatchFunc()(srcPath); err != nil {
			return err
		} else if match {
			hasMatch = true
			color.Printf("@y    CHANGED: %s\n", ev.Name)
			c.FileChanged(srcPath, ev)
		}
	}
	if !hasMatch {
		// srcPath did not match any Compiler
		if match, err := noHiddenNoIgnore(srcPath); err != nil {
			return err
		} else if match {
			color.Printf("@y    CHANGED: %s\n", ev.Name)
			fmt.Printf("Unmatched path: %s\n", srcPath)
		}
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

// recompileAllForCompiler calls RemoveOld to remove any old files the compiiler may have
// created in config.DestDir. Then it finds all the paths that match the given compiler
// (in case something changed since the last time we found the paths). Finally it compiles
// all of the matching files with a call to CompileAll.
func recompileAllForCompiler(c Compiler) error {
	if err := c.RemoveOld(); err != nil {
		return err
	}
	paths, err := FindPaths(c.CompileMatchFunc())
	if err != nil {
		return err
	}
	if err := c.CompileAll(paths); err != nil {
		return err
	}
	return nil
}

// delegateCompilePaths walks through the source directory, checks if a path matches according
// to the MatchFunc for each compiler, and adds the path to CompilerPaths if it does
// match.
func delegateCompilePaths() error {
	return filepath.Walk(config.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		matched := false
		for _, c := range Compilers {
			if match, err := c.CompileMatchFunc()(path); err != nil {
				return err
			} else if match {
				matched = true
				CompilerPaths[c] = append(CompilerPaths[c], path)
			}
		}
		if !matched && !info.IsDir() {
			// If the path didn't match any compilers according to their MatchFuncs,
			// it isn't a dir, and it is not a hidden or ignored file, add it to the
			// list of unmatched paths. These will be copied from config.SourceDir
			// to config.DestDir without being changed.
			if match, err := noHiddenNoIgnore(path); err != nil {
				return err
			} else if match {
				UnmatchedPaths = append(UnmatchedPaths, path)
			}
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
			if strings.Contains(path, string(os.PathSeparator)+".") {
				return false, nil
			}
		}
		if ignoreUnderscore {
			// Check for files and directories that begin with a '_',
			// which have special meaning in scribble and should typically
			// be ignored. If we find the substring "/_" it must mean that some
			// file or directory in the path starts with an underscore.
			if strings.Contains(path, string(os.PathSeparator)+"_") {
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
			if strings.Contains(path, string(os.PathSeparator)+".") {
				return false, nil
			}
		}
		if ignoreUnderscore {
			// Check for files and directories that begin with a '_',
			// which have special meaning in scribble and should typically
			// be ignored. If we find the substring "/_" it must mean that some
			// file or directory in the path starts with an underscore.
			// TODO: Make this compatible with windows.
			if strings.Contains(path, string(os.PathSeparator)+"_") {
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

// intersectMatchFuncs returns a MatchFunc which is functionally equivalent to the
// intersection of each MatchFunc in funcs. That is, it returns true iff each and
// every MatchFunc in funcs returns true.
func intersectMatchFuncs(funcs ...MatchFunc) MatchFunc {
	return func(path string) (bool, error) {
		for _, f := range funcs {
			if match, err := f(path); err != nil {
				return false, err
			} else if !match {
				return false, nil
			}
		}
		return true, nil
	}
}

// unionMatchFuncs returns a MatchFunc which is functionally equivalent to the
// union of each MatchFunc in funcs. That is, it returns true iff at least one MatchFunc
// in funcs returns true.
func unionMatchFuncs(funcs ...MatchFunc) MatchFunc {
	return func(path string) (bool, error) {
		for _, f := range funcs {
			if match, err := f(path); err != nil {
				return false, err
			} else if match {
				return true, nil
			}
		}
		return false, nil
	}
}

// excludeMatchFuncs returns a MatchFunc which returns true iff f returns true
// and no function in excludes returns true. It allows you to match with a simple
// function f but exclude the path if it matches some other pattern.
func excludeMatchFuncs(f MatchFunc, excludes ...MatchFunc) MatchFunc {
	return func(path string) (bool, error) {
		if firstMatch, err := f(path); err != nil {
			return false, err
		} else if !firstMatch {
			// If it doesn't match f, always return false
			return false, nil
		}
		// If it does match f, check each MatchFunc in excludes
		for _, exclude := range excludes {
			if excludeMatch, err := exclude(path); err != nil {
				return false, err
			} else if excludeMatch {
				// if path matches any MatchFunc in excludes, we should
				// exclude it. i.e., return false
				return false, nil
			}
		}
		// For all other cases, return true
		return true, nil
	}
}
