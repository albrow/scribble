package compilers

import (
	"fmt"
	"github.com/howeyc/fsnotify"
)

// HtmlTemplatesCompilerType represents a type capable of compiling go html template files.
type HtmlTemplatesCompilerType struct{}

// HtmlTemplatesCompiler is an instatiation of HtmlTemplatesCompilerType
var HtmlTemplatesCompiler = HtmlTemplatesCompilerType{}

// CompileMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".tmpl", excluding hidden and ignored
// files and directories.
func (s HtmlTemplatesCompilerType) CompileMatchFunc() MatchFunc {
	return filenameMatchFunc("*.tmpl", true, true)
}

// WatchMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".tmpl", excluding hidden files and directories,
// but including those that start with an underscore, since they may
// be imported in other files.
func (s HtmlTemplatesCompilerType) WatchMatchFunc() MatchFunc {
	return filenameMatchFunc("*.tmpl", true, false)
}

// Compile compiles the file at srcPath. The caller will only
// call this function for files which belong to HtmlTemplatesCompiler
// according to the MatchFunc. Behavior for any other file is
// undefined. Compile will output the compiled result to the appropriate
// location in config.DestDir.
func (s HtmlTemplatesCompilerType) Compile(srcPath string) error {
	return fmt.Errorf("HtmlTemplatesCompilerType.Compile not yet implemented!")
}

// CompileAll compiles zero or more files identified by srcPaths.
// It works simply by calling Compile for each path. The caller is
// responsible for only passing in files that belong to HtmlTemplatesCompiler
// according to the MatchFunc. Behavior for any other file is undefined.
func (s HtmlTemplatesCompilerType) CompileAll(srcPaths []string) error {
	return fmt.Errorf("HtmlTemplatesCompilerType.CompileAll not yet implemented!")
}

func (s HtmlTemplatesCompilerType) FileChanged(srcPath string, ev fsnotify.FileEvent) error {
	fmt.Printf("HtmlTemplatesCompiler registering change to %s\n", srcPath)
	fmt.Printf("%+v\n", ev)
	return nil
}
