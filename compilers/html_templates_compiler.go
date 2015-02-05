package compilers

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/context"
	"github.com/albrow/scribble/log"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// HtmlTemplatesCompilerType represents a type capable of compiling go html template files.
type HtmlTemplatesCompilerType struct {
	layoutFiles []string
	// createdFiles is a slice of file paths which were created by this
	// compiler. It is important for implementing the RemoveOld method.
	createdFiles []string
}

// HtmlTemplatesCompiler is an instatiation of HtmlTemplatesCompilerType
var HtmlTemplatesCompiler = HtmlTemplatesCompilerType{}

// CompileMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".tmpl", excluding hidden and ignored
// files and directories.
func (c *HtmlTemplatesCompilerType) CompileMatchFunc() MatchFunc {
	return filenameMatchFunc("*.tmpl", true, true)
}

// WatchMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".tmpl", excluding hidden files and directories,
// but including those that start with an underscore, since they may
// be imported in other files.
func (c *HtmlTemplatesCompilerType) WatchMatchFunc() MatchFunc {
	// HtmlTemplatesCompiler should watch all *tmpl files except for
	// those which are in the postsLayout dir. When those are changed,
	// they only affect posts, so we don't need to recompile any other
	// html template files.
	htmlTemplatesMatch := filenameMatchFunc("*.tmpl", true, false)
	postLayoutsMatch := pathMatchFunc(filepath.Join(config.PostLayoutsDir, "*.tmpl"), true, false)
	// excludeMatchFuncs lets us express these conditions easily. It
	// returns a MatchFunc which will return true iff the path represents
	// and html template *and* is *not* in the post layouts dir. I.e., if
	// a .tmpl file is in the post layouts dir, it will return false and
	// HtmlTemplatesCompiler will not be alerted when those files change.
	return excludeMatchFuncs(htmlTemplatesMatch, postLayoutsMatch)
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
func (c *HtmlTemplatesCompilerType) Compile(srcPath string) error {
	// parse path and figure out destPath
	destPath := strings.Replace(srcPath, ".tmpl", ".html", 1)
	destPath = strings.Replace(destPath, config.SourceDir, config.DestDir, 1)
	log.Success.Printf("CREATE: %s -> %s", srcPath, destPath)

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

	// Create the template by parsing the raw content. Then parse all the layout files, include files, and add context.FuncMap
	tmpl := template.New(filepath.Base(srcPath))
	tmpl.Funcs(context.FuncMap)
	if _, err := tmpl.Parse(content); err != nil {
		return err
	}
	if config.LayoutsDir != "" {
		if _, err := tmpl.ParseGlob(filepath.Join(config.LayoutsDir, "*.tmpl")); err != nil {
			return err
		}
	} else {
		// config.LayoutsDir is more or less required. Every page must have a layout
		return fmt.Errorf("Missing required config variable: layoutsDir. Please add it to config.toml.")
	}
	if config.IncludesDir != "" {
		// config.IncludesDir, on the other hand, is optional. You don't have to use includes.
		if _, err := tmpl.ParseGlob(filepath.Join(config.IncludesDir, "*.tmpl")); err != nil {
			return err
		}
	}

	// Create and write to the destination file
	destFile, err := util.CreateFileWithPath(destPath)
	if err != nil {
		return err
	}
	if err := tmpl.ExecuteTemplate(destFile, layout, pageContext); err != nil {
		return err
	}

	// Add the created file to the list of created files
	c.createdFiles = append(c.createdFiles, destPath)

	return nil
}

// CompileAll compiles zero or more files identified by srcPaths.
// It works simply by calling Compile for each path. The caller is
// responsible for only passing in files that belong to HtmlTemplatesCompiler
// according to the MatchFunc. Behavior for any other file is undefined.
func (c *HtmlTemplatesCompilerType) CompileAll(srcPaths []string) error {
	log.Default.Println("Compiling go html templates...")
	for _, srcPath := range srcPaths {
		if err := c.Compile(srcPath); err != nil {
			return err
		}
	}
	return nil
}

func (c *HtmlTemplatesCompilerType) FileChanged(srcPath string, ev *fsnotify.FileEvent) error {
	// TODO: Analyze template files and be more intelligent here?
	// If a single file was changed, only recompile that file. If a
	// layout file was changed, recompile all the files that use that
	// layout. For now, just recompile all html templates.
	if err := recompileAllForCompiler(c); err != nil {
		return err
	}
	return nil
}

func (c *HtmlTemplatesCompilerType) RemoveOld() error {
	// Simply iterate through createdFiles and remove each of them
	for _, path := range c.createdFiles {
		if err := util.RemoveIfExists(path); err != nil {
			return err
		}
	}
	return nil
}
