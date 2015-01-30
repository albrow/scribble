package compilers

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/test_util"
	"github.com/albrow/scribble/util"
	"os"
	"path/filepath"
	"testing"
)

func TestSassPathMatch(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := string(os.PathSeparator) + filepath.Join("tmp", "sass_compiler_paths")
	defer func() {
		// Remove everything after we're done
		if err := os.RemoveAll(root); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}
	}()

	// Create a few files.
	tmpPaths := []string{
		filepath.Join(root, "styles", "main.scss"),
		filepath.Join(root, "styles", "_colors.scss"),
		filepath.Join(root, "styles", "_body.scss"),
		filepath.Join(root, "styles", "notice.txt"),
		filepath.Join(root, "styles", "README"),
		filepath.Join(root, "_sass", "main.scss"),
		filepath.Join(root, ".sass", "main.scss"),
		filepath.Join(root, "more_sass", "other_stuff", "this.scss"),
	}
	if err := util.CreateEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the SassCompiler,
	// the other files should be ignored.
	expectedPaths := []string{
		filepath.Join(root, "styles", "main.scss"),
		filepath.Join(root, "more_sass", "other_stuff", "this.scss"),
	}

	// Use the MatchFunc to find all the paths
	config.SourceDir = root
	gotPaths, err := FindPaths(SassCompiler.CompileMatchFunc())
	if err != nil {
		t.Error(err)
	}

	// Check that the paths we get are correct
	test_util.CheckStringsMatch(t, expectedPaths, gotPaths)
}

func TestSassCompile(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := string(os.PathSeparator) + filepath.Join("tmp", "test_sass_compiler")
	defer func() {
		// Remove everything after we're done
		if err := os.RemoveAll(root); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}
	}()

	// Copy some files from test_files to source directory in the temp root
	testFilesDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "albrow", "scribble", "test_files", "sass")
	srcDir := filepath.Join(root, "source")
	destDir := filepath.Join(root, "public")
	if err := util.RecursiveCopy(testFilesDir+"/source", srcDir); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the sass files
	config.SourceDir = filepath.Join(root, "source")
	config.DestDir = filepath.Join(root, "public")
	if err := SassCompiler.Compile(srcDir + "/styles/main.scss"); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	expectedDir := filepath.Join(testFilesDir, "public")
	expectedFile := filepath.Join(expectedDir, "styles", "main.css")
	gotFile := filepath.Join(destDir, "styles", "main.css")
	test_util.CheckFilesMatch(t, expectedFile, gotFile)
}
