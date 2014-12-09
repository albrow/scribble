package generators

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/util"
	"os"
	"testing"
)

func TestSassPathMatch(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/sass_compiler"
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
		root + "/styles/main.scss",
		root + "/styles/_colors.scss",
		root + "/styles/_body.scss",
		root + "/styles/notice.txt",
		root + "/styles/README",
		root + "/_sass/main.scss",
		root + "/.sass/main.scss",
		root + "/more_sass/other_stuff/this.scss",
	}
	if err := createEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the SassCompiler,
	// the other files should be ignored.
	expectedPaths := []string{
		root + "/styles/main.scss",
		root + "/more_sass/other_stuff/this.scss",
	}

	// Check that the paths we get are correct
	checkPathsMatch(t, SassCompiler, root, expectedPaths)
}

func TestSassCompile(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/test_sass_compiler"
	defer func() {
		// Remove everything after we're done
		if err := os.RemoveAll(root); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}
	}()

	// Copy some files from test_files to source directory in the temp root
	testFilesDir := os.Getenv("GOPATH") + "/src/github.com/albrow/scribble/test_files/sass"
	srcDir := root + "/source"
	destDir := root + "/public"
	if err := util.RecursiveCopy(testFilesDir+"/source", srcDir); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the sass files
	config.SourceDir = root + "/source"
	config.DestDir = root + "/public"
	if err := SassCompiler.Compile(srcDir + "/styles/main.scss"); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	expectedDir := testFilesDir + "/public"
	checkOutputMatchesFile(t, destDir+"/styles/main.css", expectedDir+"/styles/main.css")
}
