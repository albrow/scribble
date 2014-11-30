package generators

import (
	"fmt"
	"github.com/albrow/scribble/lib"
	"github.com/wsxiaoys/terminal/color"
	"path/filepath"
	"strings"
)

type PostsCompilerType struct {
	pathMatch string
}

var PostsCompiler = PostsCompilerType{
	pathMatch: "",
}

func (p *PostsCompilerType) Init() {
	p.pathMatch = lib.PostsDir + "/*.md"

}

func (p PostsCompilerType) GetWalkFunc(paths *[]string) filepath.WalkFunc {
	return pathMatchWalkFunc(paths, p.pathMatch, true, false)
}

func (p PostsCompilerType) Compile(srcPath string, destDir string) error {
	// Get the parsed post object and determine dest path
	post := lib.GetOrCreatePostFromPath(srcPath)
	srcFilename := filepath.Base(srcPath)
	destPath := fmt.Sprintf("%s/%s", destDir, strings.TrimSuffix(srcFilename, ".md"))
	destIndexFilePath := destPath + "/index.html"
	color.Printf("@g    CREATE: %s -> %s\n", srcPath, destIndexFilePath)

	// Create the index file
	fmt.Println("creating file")
	destFile, err := lib.CreateFileWithPath(destIndexFilePath)
	if err != nil {
		return err
	}

	// Get and compile the template
	fmt.Println("getting template")
	tmpl := lib.GetPostTemplate()
	fmt.Println("parsing post")
	post.Parse()
	fmt.Println("getting context")
	postContext := lib.GetContext()
	postContext["Post"] = post
	fmt.Println("Executing template")
	if err := tmpl.Execute(destFile, postContext); err != nil {
		return fmt.Errorf("ERROR compiling ace template for posts: %s", err.Error())
	}
	fmt.Println("Done")
	return nil
}

func (p PostsCompilerType) CompileAll(srcPaths []string, destDir string) error {
	return nil
}
