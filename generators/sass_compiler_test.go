package generators

import (
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
		root + "sass/main.scss",
		root + "sass/_colors.scss",
		root + "sass/_body.scss",
		root + "sass/notice.txt",
		root + "sass/README",
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
		root + "sass/main.scss",
		root + "more_sass/other_stuff/this.scss",
	}

	// Check that the paths we get are correct
	checkPathsMatch(t, SassCompiler, root, expectedPaths)

	// Copy some files from test_files to the tmp directory
	gopath := os.Getenv("GOPATH")
	srcRoot := gopath + "/src/github.com/albrow/scribble/test_files/sass/"
	pathsToCopy := map[string]string{
		root + "sass/main.scss":    srcRoot + "main.scss",
		root + "sass/_colors.scss": srcRoot + "_colors.scss",
		root + "sass/_body.scss":   srcRoot + "_body.scss",
	}
	if err := copyFiles(pathsToCopy); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the sass files
	if err := SassCompiler.Compile(root+"sass/main.scss", root+"public"); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	checkOutputMatchesFile(t, root+"public/main.css", srcRoot+"expected.css")
}
