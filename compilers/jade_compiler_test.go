// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package compilers

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/test_util"
	"github.com/albrow/scribble/util"
	"os"
	"path/filepath"
	"testing"
)

func TestJadePathMatch(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := string(os.PathSeparator) + filepath.Join("tmp", "jade_compiler_paths")
	defer func() {
		// Remove everything after we're done
		if err := util.RemoveAllIfExists(root); err != nil {
			panic(err)
		}
	}()

	// Create a few files.
	tmpPaths := []string{
		filepath.Join(root, "index.jade"),
		filepath.Join(root, "archive", "index.jade"),
		filepath.Join(root, "_includes", "foot.jade"),
		filepath.Join(root, "_templates", "base.jade"),
		filepath.Join(root, "archive", "notice.txt"),
		filepath.Join(root, ".cache", "index.jade"),
	}
	if err := util.CreateEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the JadeCompiler,
	// the other files should be ignored.
	expectedPaths := []string{
		filepath.Join(root, "index.jade"),
		filepath.Join(root, "archive", "index.jade"),
	}

	// Use the MatchFunc to find all the paths
	config.SourceDir = root
	gotPaths, err := FindPaths(JadeCompiler.CompileMatchFunc())
	if err != nil {
		t.Error(err)
	}

	// Check that the paths we get are correct
	test_util.CheckStringsMatch(t, expectedPaths, gotPaths)
}

func TestJadeCompile(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := string(os.PathSeparator) + filepath.Join("tmp", "test_jade_compiler")
	defer func() {
		// Remove everything after we're done
		if err := util.RemoveAllIfExists(root); err != nil {
			panic(err)
		}
	}()

	// Copy some files from test_files to source directory in the temp root
	testFilesDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "albrow", "scribble", "test_files", "jade")
	srcDir := filepath.Join(root, "source")
	destDir := filepath.Join(root, "public")
	if err := util.RecursiveCopy(filepath.Join(testFilesDir, "source"), srcDir); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the html template files
	config.SourceDir = filepath.Join(root, "source")
	config.DestDir = filepath.Join(root, "public")
	if err := JadeCompiler.Compile(filepath.Join(srcDir, "index.jade")); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	expectedDir := filepath.Join(testFilesDir, "public")
	test_util.CheckFilesMatch(t, filepath.Join(expectedDir, "index.html"), filepath.Join(destDir, "index.html"))
}
