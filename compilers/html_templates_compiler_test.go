package compilers

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/test_util"
	"github.com/albrow/scribble/util"
	"os"
	"path/filepath"
	"testing"
)

func TestHtmlTemplatesPathMatch(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := string(os.PathSeparator) + filepath.Join("tmp", "html_templates_compiler_paths")
	defer func() {
		// Remove everything after we're done
		if err := util.RemoveAllIfExists(root); err != nil {
			panic(err)
		}
	}()

	// Create a few files.
	tmpPaths := []string{
		filepath.Join(root, "index.tmpl"),
		filepath.Join(root, "pages", "about.tmpl"),
		filepath.Join(root, "pages", "_partial.tmpl"),
		filepath.Join(root, "notice.txt"),
		filepath.Join(root, "other", "README"),
		filepath.Join(root, "_layouts", "main.tmpl"),
		filepath.Join(root, ".build", "tmpl", "main.tmpl"),
		filepath.Join(root, "more_pages", "other_stuff", "page.tmpl"),
	}
	if err := util.CreateEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the HtmlTemplatesCompiler,
	// the other files should be ignored.
	expectedPaths := []string{
		filepath.Join(root, "index.tmpl"),
		filepath.Join(root, "pages", "about.tmpl"),
		filepath.Join(root, "more_pages", "other_stuff", "page.tmpl"),
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

func TestHtmlTemplatesCompile(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := string(os.PathSeparator) + filepath.Join("tmp", "test_html_templates_compiler")
	defer func() {
		// Remove everything after we're done
		if err := util.RemoveAllIfExists(root); err != nil {
			panic(err)
		}
	}()

	// Copy some files from test_files to source directory in the temp root
	testFilesDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "albrow", "scribble", "test_files", "html_templates")
	srcDir := filepath.Join(root, "source")
	destDir := filepath.Join(root, "public")
	if err := util.RecursiveCopy(filepath.Join(testFilesDir, "source"), srcDir); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the html template files
	config.SourceDir = filepath.Join(root, "source")
	config.DestDir = filepath.Join(root, "public")
	config.LayoutsDir = filepath.Join(config.SourceDir, "_layouts")
	config.IncludesDir = filepath.Join(config.SourceDir, "_includes")
	if err := HtmlTemplatesCompiler.Compile(filepath.Join(srcDir, "index.tmpl")); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	expectedDir := filepath.Join(testFilesDir, "public")
	test_util.CheckFilesMatch(t, filepath.Join(expectedDir, "index.html"), filepath.Join(destDir, "index.html"))
}
