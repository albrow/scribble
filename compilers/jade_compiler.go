package compilers

import (
	"encoding/json"
	"fmt"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/context"
	"github.com/albrow/scribble/log"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// JadeCompilerType represents a type capable of compiling jade files.
type JadeCompilerType struct {
	// createdFiles is a slice of file paths which were created by this
	// compiler. It is important for implementing the RemoveOld method.
	createdFiles []string
}

// JadeCompiler is an instatiation of JadeCompilerType
var JadeCompiler = JadeCompilerType{}

// CompileMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".jade", excluding hidden and ignored
// files and directories.
func (*JadeCompilerType) CompileMatchFunc() MatchFunc {
	return filenameMatchFunc("*.jade", true, true)
}

// WatchMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that ends in ".jade", excluding hidden files and directories,
// but including those that start with an underscore, since they may
// be imported in other files.
func (*JadeCompilerType) WatchMatchFunc() MatchFunc {
	return filenameMatchFunc("*.jade", true, false)
}

// Compile compiles the file at srcPath. The caller will only
// call this function for files which belong to JadeCompiler
// according to the MatchFunc. Behavior for any other file is
// undefined. Compile will output the compiled result to the appropriate
// location in config.DestDir.
func (j *JadeCompilerType) Compile(srcPath string) error {
	// parse path and figure out destPath
	destPath := strings.Replace(srcPath, ".jade", ".html", 1)
	destPath = strings.Replace(destPath, config.SourceDir, config.DestDir, 1)
	log.Success.Printf("CREATE: %s -> %s", srcPath, destPath)

	// create the dest directory if needed
	if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
		// if the dir already exists, that's fine
		// if there was some other error, return it
		if !os.IsExist(err) {
			return err
		}
	}

	// set up and execute the command, capturing the output only if there was an error
	destDir := filepath.Dir(destPath)
	cmd := exec.Command("jade", srcPath, "--out", destDir)
	response, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("while compiling jade: %s", string(response))
	}

	// Add destPath to the list of created files
	j.createdFiles = append(j.createdFiles, destPath)

	return nil
}

// CompileAll compiles zero or more files identified by srcPaths.
// It works simply by calling Compile for each path. The caller is
// responsible for only passing in files that belong to JadeCompiler
// according to the MatchFunc. Behavior for any other file is undefined.
func (j *JadeCompilerType) CompileAll(srcPaths []string) error {
	log.Default.Println("Compiling jade...")
	for _, srcPath := range srcPaths {
		if err := j.Compile(srcPath); err != nil {
			return err
		}
	}
	return nil
}

func (j *JadeCompilerType) FileChanged(srcPath string, ev *fsnotify.FileEvent) error {
	// TODO: Analyze jade files and be more intelligent here?
	// Only recompile the file at srcPath and any files that import it?
	// For now, just recompile all jade.
	if err := recompileAllForCompiler(j); err != nil {
		return err
	}
	return nil
}

func (j *JadeCompilerType) RemoveOld() error {
	// Simply iterate through createdFiles and remove each of them
	for _, path := range j.createdFiles {
		if err := util.RemoveIfExists(path); err != nil {
			return err
		}
	}
	return nil
}

func (c *JadeCompilerType) PostLayoutMatchFunc() MatchFunc {
	return filenameMatchFunc("*.jade", true, false)
}

func (c *JadeCompilerType) RenderPost(post *Post, destPath string) error {
	// Create the context for the post (in this case json data)
	postContext := context.CopyContext()
	postContext["Post"] = post
	jsonContext, err := json.Marshal(postContext)
	if err != nil {
		return fmt.Errorf("ERROR converting post context to json for jade post layout:\n%s", err.Error())
	}

	// set up and execute the command, capturing the output only if there was an error
	postLayoutFile := filepath.Join(config.PostLayoutsDir, post.LayoutName)
	destDir := filepath.Dir(destPath)
	cmd := exec.Command("jade", postLayoutFile, "--out", destDir, "--obj", string(jsonContext))
	response, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("while compiling jade: %s", string(response))
	}

	// jade does not allow us to specify the filename, so we'll manually do a rename
	// TODO: on unixy systems use a pipe or redirect to a file
	layoutNameExt := filepath.Ext(post.LayoutName)
	layoutNameNoExt := strings.Replace(post.LayoutName, layoutNameExt, "", 1)
	oldName := strings.Replace(destPath, "index", layoutNameNoExt, 1)
	if err := os.Rename(oldName, destPath); err != nil {
		return err
	}
	return nil
}
