package generators

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/ace"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/context"
	"github.com/albrow/scribble/util"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"path/filepath"
	"strings"
)

type AceCompilerType struct {
	filenameMatch string
}

const aceFilenameMatch = "*.ace"

var AceCompiler = AceCompilerType{
	filenameMatch: aceFilenameMatch,
}

func (a AceCompilerType) GetWalkFunc(paths *[]string) filepath.WalkFunc {
	return filenameMatchWalkFunc(paths, a.filenameMatch, true, true)
}

func (a AceCompilerType) Compile(srcPath string, destDir string) error {
	// parse path and figure out destPath
	destPath := strings.Replace(srcPath, ".ace", ".html", 1)
	destPath = strings.Replace(destPath, config.SourceDir, destDir, 1)
	color.Printf("@g    CREATE: %s -> %s\n", srcPath, destPath)

	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		panic(err)
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

func (a AceCompilerType) CompileAll(srcPaths []string, destDir string) error {
	fmt.Println("--> compiling ace")
	for _, srcPath := range srcPaths {
		if err := a.Compile(srcPath, destDir); err != nil {
			return err
		}
	}
	return nil
}
