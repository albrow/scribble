package compilers

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/test_util"
	"github.com/albrow/scribble/util"
	"os"
	"testing"
)

func TestHtmlTemplatesPathMatch(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/html_templates_compiler"
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
		root + "/index.tmpl",
		root + "/pages/about.tmpl",
		root + "/pages/_partial.tmpl",
		root + "/notice.txt",
		root + "/other/README",
		root + "/_layouts/main.tmpl",
		root + "/.build/tmpl/main.tmpl",
		root + "/more_pages/other_stuff/page.tmpl",
	}
	if err := util.CreateEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the HtmlTemplatesCompiler,
	// the other files should be ignored.
	expectedPaths := []string{
		root + "/index.tmpl",
		root + "/pages/about.tmpl",
		root + "/more_pages/other_stuff/page.tmpl",
	}

	// Use the MatchFunc to find all the paths
	config.SourceDir = root
	gotPaths, err := FindPaths(HtmlTemplatesCompiler.CompileMatchFunc())
	if err != nil {
		t.Error(err)
	}

	// Check that the paths we get are correct
	test_util.CheckStringsMatch(t, expectedPaths, gotPaths)
}
