package compilers

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/ace"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/context"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"strings"
)

// AceCompilerType represents a type capable of compiling ace templates.
type AceCompilerType struct{}

// AceCompiler is an instatiation of AceCompilerType
var AceCompiler = AceCompilerType{}

// CompileMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".ace", excluding hidden and ignored
// files and directories.
func (a AceCompilerType) CompileMatchFunc() MatchFunc {
	return filenameMatchFunc("*.ace", true, true)
}

// WatchMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".ace", excluding hidden files and directories
// but including those that start with an underscore. Files which
// start with an underscore may be included in other files, so we
// need to watch them too.
func (a AceCompilerType) WatchMatchFunc() MatchFunc {
	return filenameMatchFunc("*.ace", true, false)
}

// Compile compiles the file at srcPath. The caller will only
// call this function for files which belong to AceCompiler
// according to the MatchFunc. Behavior for any other file is
// undefined. Compile will output the compiled result to the appropriate
// location in config.DestDir.
func (a AceCompilerType) Compile(srcPath string) error {
	// parse path and figure out destPath
	destPath := strings.Replace(srcPath, ".ace", ".html", 1)
	destPath = strings.Replace(destPath, config.SourceDir, config.DestDir, 1)
	color.Printf("@g    CREATE: %s -> %s\n", srcPath, destPath)

	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(srcFile)

	// Split source file into front matter and content
	frontMatter, content, err := util.SplitFrontMatter(reader)
	pageContext := context.GetContext()
	if frontMatter != "" {
		if _, err := toml.Decode(frontMatter, pageContext); err != nil {
			return err
		}
	}

	// Determine the correct layout and render the template
	layout := "base"
	if otherLayout, found := pageContext["layout"]; found {
		layout = otherLayout.(string)
	}
	layoutPath := config.LayoutsDir + "/" + layout
	tpl, err := ace.Load(layoutPath, srcPath, &ace.Options{
		DynamicReload: true,
		BaseDir:       config.SourceDir,
		FuncMap:       context.FuncMap,
		Asset: func(name string) ([]byte, error) {
			return []byte(content), nil
		},
	})
	if err != nil {
		return err
	}

	destFile, err := util.CreateFileWithPath(destPath)
	if err != nil {
		return err
	}
	if err := tpl.Execute(destFile, pageContext); err != nil {
		return err
	}
	return nil
}

// CompileAll compiles zero or more files identified by srcPaths.
// It works simply by calling Compile for each path. The caller is
// responsible for only passing in files that belong to AceCompiler
// according to the MatchFunc. Behavior for any other file is undefined.
func (a AceCompilerType) CompileAll(srcPaths []string) error {
	fmt.Println("--> compiling ace")
	for _, srcPath := range srcPaths {
		if err := a.Compile(srcPath); err != nil {
			return err
		}
	}
	return nil
}

func (a AceCompilerType) FileChanged(srcPath string, ev fsnotify.FileEvent) error {
	fmt.Printf("AceCompiler registering change to %s\n", srcPath)
	fmt.Printf("%+v\n", ev)
	return nil
}
