package generators

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/ace"
	"github.com/albrow/scribble/lib"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"path/filepath"
	"strings"
)

type AceCompilerType struct {
	pathMatch string
}

const acePathMatch = "*.ace"

var AceCompiler = AceCompilerType{
	pathMatch: acePathMatch,
}

func (a AceCompilerType) PathMatch() string {
	return a.pathMatch
}

func (a AceCompilerType) IgnoreHidden() bool {
	return true
}

func (a AceCompilerType) IgnoreUnderscore() bool {
	return true
}

func (a AceCompilerType) Compile(srcPath string, destDir string) error {
	// Parse path and figure out destPath
	srcFilename := filepath.Base(srcPath)
	destFilename := strings.Replace(srcFilename, ".ace", ".html", 1)
	destPath := fmt.Sprintf("%s/%s", destDir, destFilename)
	color.Printf("@g    CREATE: %s -> %s\n", srcPath, destPath)

	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(srcFile)

	// Split source file into front matter and content
	frontMatter, content, err := lib.SplitFrontMatter(reader)
	pageContext := lib.GetContext()
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
	layoutPath := lib.LayoutsDir + "/" + layout
	tpl, err := ace.Load(layoutPath, srcPath, &ace.Options{
		DynamicReload: true,
		BaseDir:       lib.SourceDir,
		FuncMap:       lib.FuncMap,
		Asset: func(name string) ([]byte, error) {
			return []byte(content), nil
		},
	})
	if err != nil {
		return err
	}

	destFile, err := lib.CreateFileWithPath(destPath)
	if err != nil {
		return err
	}
	if err := tpl.Execute(destFile, pageContext); err != nil {
		return err
	}
	return nil
}

func (a AceCompilerType) CompileAll(srcPaths []string, destDir string) error {
	fmt.Println("Call to compile all: ", srcPaths)
	return nil
}
