package compilers

import (
	"fmt"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SassCompilerType represents a type capable of compiling sass files.
type SassCompilerType struct {
	// createdFiles is a slice of file paths which were created by this
	// compiler. It is important for implementing the RemoveOld method.
	createdFiles []string
}

// SassCompiler is an instatiation of SassCompilerType
var SassCompiler = SassCompilerType{}

// CompileMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".scss", excluding hidden and ignored
// files and directories.
func (s *SassCompilerType) CompileMatchFunc() MatchFunc {
	return filenameMatchFunc("*.scss", true, true)
}

// WatchMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".scss", excluding hidden files and directories,
// but including those that start with an underscore, since they may
// be imported in other files.
func (s *SassCompilerType) WatchMatchFunc() MatchFunc {
	return filenameMatchFunc("*.scss", true, false)
}

// Compile compiles the file at srcPath. The caller will only
// call this function for files which belong to SassCompiler
// according to the MatchFunc. Behavior for any other file is
// undefined. Compile will output the compiled result to the appropriate
// location in config.DestDir.
func (s *SassCompilerType) Compile(srcPath string) error {
	// parse path and figure out destPath
	destPath := strings.Replace(srcPath, ".scss", ".css", 1)
	destPath = strings.Replace(destPath, config.SourceDir, config.DestDir, 1)
	color.Printf("@g    CREATE: %s -> %s\n", srcPath, destPath)

	// create the dest directory if needed
	if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
		// if the dir already exists, that's fine
		// if there was some other error, return it
		if !os.IsExist(err) {
			return err
		}
	}

	// set up and execute the command, capturing the output only if there was an error
	cmd := exec.Command("sassc", srcPath, destPath)
	response, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ERROR compiling sass: %s", string(response))
	}

	// Add destPath to the list of created files
	s.createdFiles = append(s.createdFiles, destPath)

	return nil
}

// CompileAll compiles zero or more files identified by srcPaths.
// It works simply by calling Compile for each path. The caller is
// responsible for only passing in files that belong to SassCompiler
// according to the MatchFunc. Behavior for any other file is undefined.
func (s *SassCompilerType) CompileAll(srcPaths []string) error {
	fmt.Println("--> compiling sass")
	for _, srcPath := range srcPaths {
		if err := s.Compile(srcPath); err != nil {
			return err
		}
	}
	return nil
}

func (s *SassCompilerType) FileChanged(srcPath string, ev fsnotify.FileEvent) error {
	// TODO: Analyze sass files and be more intelligent here?
	// Only recompile the file at srcPath and any files that import it?
	// For now, just recompile all sass.
	if err := recompileAllForCompiler(s); err != nil {
		return err
	}
	return nil
}

func (s *SassCompilerType) RemoveOld() error {
	// Simply iterate through createdFiles and remove each of them
	for _, path := range s.createdFiles {
		if err := util.RemoveIfExists(path); err != nil {
			return err
		}
	}
	return nil
}
