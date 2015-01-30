package compilers

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/context"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"github.com/wsxiaoys/terminal/color"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// HtmlTemplatesCompilerType represents a type capable of compiling go html template files.
type HtmlTemplatesCompilerType struct {
	layoutFiles []string
}

// HtmlTemplatesCompiler is an instatiation of HtmlTemplatesCompilerType
var HtmlTemplatesCompiler = HtmlTemplatesCompilerType{}

// CompileMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".tmpl", excluding hidden and ignored
// files and directories.
func (c HtmlTemplatesCompilerType) CompileMatchFunc() MatchFunc {
	return filenameMatchFunc("*.tmpl", true, true)
}

// WatchMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".tmpl", excluding hidden files and directories,
// but including those that start with an underscore, since they may
// be imported in other files.
func (c HtmlTemplatesCompilerType) WatchMatchFunc() MatchFunc {
	return filenameMatchFunc("*.tmpl", true, false)
}

// Init should be called before any other methods. In this case, Init
// finds and loads the layout templates in config.LayoutsDir
func (c *HtmlTemplatesCompilerType) Init() {
	pattern := filepath.Join(config.LayoutsDir, "*.tmpl")
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	c.layoutFiles = files
}

// Compile compiles the file at srcPath. The caller will only
// call this function for files which belong to HtmlTemplatesCompiler
// according to the MatchFunc. Behavior for any other file is
// undefined. Compile will output the compiled result to the appropriate
// location in config.DestDir.
func (c HtmlTemplatesCompilerType) Compile(srcPath string) error {
	// parse path and figure out destPath
	destPath := strings.Replace(srcPath, ".tmpl", ".html", 1)
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
	pageContext := context.CopyContext()
	if frontMatter != "" {
		if _, err := toml.Decode(frontMatter, pageContext); err != nil {
			return err
		}
	}

	// Read the layout key from the toml frontmatter
	layoutKey, found := pageContext["layout"]
	if !found {
		return fmt.Errorf("Could not find layout definition in toml frontmatter for html template: %s", srcPath)
	}
	layout, ok := layoutKey.(string)
	if !ok {
		return fmt.Errorf("Could not convert frontmatter key layout of type %T to string!", layoutKey)
	}

	// Create the template by parsing the raw content. Then parse all the layout files and add context.FuncMap
	tmpl := template.New(filepath.Base(srcPath))
	tmpl.Funcs(context.FuncMap)
	if _, err := tmpl.Parse(content); err != nil {
		return err
	}
	if _, err := tmpl.ParseGlob(filepath.Join(config.LayoutsDir, "*.tmpl")); err != nil {
		return err
	}

	// Create and write to the destination file
	destFile, err := util.CreateFileWithPath(destPath)
	if err != nil {
		return err
	}
	if err := tmpl.ExecuteTemplate(destFile, layout, pageContext); err != nil {
		return err
	}

	return nil
}

// CompileAll compiles zero or more files identified by srcPaths.
// It works simply by calling Compile for each path. The caller is
// responsible for only passing in files that belong to HtmlTemplatesCompiler
// according to the MatchFunc. Behavior for any other file is undefined.
func (c HtmlTemplatesCompilerType) CompileAll(srcPaths []string) error {
	fmt.Println("--> compiling go html templates...")
	for _, srcPath := range srcPaths {
		if err := c.Compile(srcPath); err != nil {
			return err
		}
	}
	return nil
}

func (c HtmlTemplatesCompilerType) FileChanged(srcPath string, ev fsnotify.FileEvent) error {
	fmt.Printf("HtmlTemplatesCompiler registering change to %s\n", srcPath)
	fmt.Printf("%+v\n", ev)
	return nil
}
