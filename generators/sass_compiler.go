package generators

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SassCompilerType struct {
	pathMatch string
}

const sassPathMatch = "*.scss"

var SassCompiler = SassCompilerType{
	pathMatch: sassPathMatch,
}

func (s SassCompilerType) PathMatch() string {
	return s.pathMatch
}

func (s SassCompilerType) IgnoreHidden() bool {
	return true
}

func (s SassCompilerType) IgnoreUnderscore() bool {
	return true
}

func (s SassCompilerType) Compile(srcPath string, destDir string) error {
	// parse path and figure out destPath
	srcFilename := filepath.Base(srcPath)
	destFilename := strings.Replace(srcFilename, ".scss", ".css", 1)
	destPath := fmt.Sprintf("%s/%s", destDir, destFilename)
	color.Printf("@g    CREATE: %s -> %s\n", srcPath, destPath)

	// create the destDir if needed
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
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
	return nil
}

func (s SassCompilerType) CompileAll(srcPaths []string, destDir string) error {
	fmt.Println("Call to compile all: ", srcPaths)
	return nil
}
