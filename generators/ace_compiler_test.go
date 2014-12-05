package generators

import (
	"github.com/albrow/scribble/lib"
	"os"
	"testing"
)

func TestAceCompile(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/ace_compiler/"
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
		root + "index.ace",
		root + "_layouts/base.ace",
		root + "_includes/greet.ace",
		root + "ace/notice.txt",
		root + "ace/README",
		root + ".templates/base.ace",
		root + "more/other_stuff/this.ace",
	}
	if err := createEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the AceCompiler,
	// the other files should be ignored.
	expectedPaths := []string{
		root + "index.ace",
		root + "more/other_stuff/this.ace",
	}

	// Check that the paths we get are correct
	checkPathsMatch(t, AceCompiler, root, expectedPaths)

	// Copy some files from test_files to the tmp directory
	gopath := os.Getenv("GOPATH")
	srcRoot := gopath + "/src/github.com/albrow/scribble/test_files/ace/"
	pathsToCopy := map[string]string{
		root + "_layouts/base.ace":   srcRoot + "_layouts/base.ace",
		root + "index.ace":           srcRoot + "index.ace",
		root + "_includes/greet.ace": srcRoot + "_includes/greet.ace",
	}
	if err := copyFiles(pathsToCopy); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the ace files
	lib.LayoutsDir = "_layouts"
	lib.SourceDir = root
	if err := AceCompiler.Compile(root+"index.ace", root+"public"); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	checkOutputMatchesFile(t, root+"public/index.html", srcRoot+"expected.html")
}
