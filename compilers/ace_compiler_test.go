package compilers

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/test_util"
	"github.com/albrow/scribble/util"
	"os"
	"testing"
)

func TestAcePathMatch(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/test_ace_compiler_paths"
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
		root + "/index.ace",
		root + "/_layouts/base.ace",
		root + "/_includes/greet.ace",
		root + "/ace/notice.txt",
		root + "/ace/README",
		root + "/.templates/base.ace",
		root + "/more/other_stuff/this.ace",
	}
	if err := util.CreateEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the AceCompiler,
	// the other files should be ignored.
	expectedPaths := []string{
		root + "/index.ace",
		root + "/more/other_stuff/this.ace",
	}

	// Use the MatchFunc to find all the paths
	config.SourceDir = root
	gotPaths, err := FindPaths(AceCompiler.CompileMatchFunc())
	if err != nil {
		t.Error(err)
	}

	// Check that the paths we get are correct
	test_util.CheckStringsMatch(t, expectedPaths, gotPaths)
}

func TestAceCompile(t *testing.T) {
	t.Skip("Skipping Ace Compiler test. Going to switch to std html/template soon.")
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/test_ace_compiler"
	defer func() {
		// Remove everything after we're done
		if err := os.RemoveAll(root); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}
	}()

	// Copy some files from test_files to source directory in the temp root
	testFilesDir := os.Getenv("GOPATH") + "/src/github.com/albrow/scribble/test_files/ace"
	srcDir := root + "/source"
	destDir := root + "/public"
	if err := util.RecursiveCopy(testFilesDir+"/source", srcDir); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the ace files
	config.LayoutsDir = "_layouts"
	config.SourceDir = srcDir
	config.DestDir = destDir
	if err := AceCompiler.Compile(srcDir + "/index.ace"); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	expectedDir := testFilesDir + "/public"
	test_util.CheckFilesMatch(t, expectedDir+"/index.html", destDir+"/index.html")
}
