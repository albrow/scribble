package generators

import (
	"github.com/albrow/scribble/config"
	"os"
	"testing"
)

func TestSassCompile(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/sass_compiler/"
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
		root + "styles/main.scss",
		root + "styles/_colors.scss",
		root + "styles/_body.scss",
		root + "styles/notice.txt",
		root + "styles/README",
		root + "_sass/main.scss",
		root + ".sass/main.scss",
		root + "more_sass/other_stuff/this.scss",
	}
	if err := createEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the SassCompiler,
	// the other files should be ignored.
	expectedPaths := []string{
		root + "styles/main.scss",
		root + "more_sass/other_stuff/this.scss",
	}

	// Check that the paths we get are correct
	checkPathsMatch(t, SassCompiler, root, expectedPaths)

	// Copy some files from test_files to the tmp directory
	gopath := os.Getenv("GOPATH")
	srcRoot := gopath + "/src/github.com/albrow/scribble/test_files/sass/"
	pathsToCopy := map[string]string{
		root + "styles/main.scss":    srcRoot + "main.scss",
		root + "styles/_colors.scss": srcRoot + "_colors.scss",
		root + "styles/_body.scss":   srcRoot + "_body.scss",
	}
	if err := copyFiles(pathsToCopy); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the sass files
	config.SourceDir = root[0 : len(root)-1]
	if err := SassCompiler.Compile(root+"styles/main.scss", root+"public"); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	checkOutputMatchesFile(t, root+"public/styles/main.css", srcRoot+"expected.css")
}
