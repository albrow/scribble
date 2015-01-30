package compilers

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/test_util"
	"github.com/albrow/scribble/util"
	"os"
	"path/filepath"
	"testing"
)

func TestPostsPathMatch(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := string(os.PathSeparator) + filepath.Join("tmp", "posts_compiler_paths")
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
		filepath.Join(root, "_posts", "post.md"),
		filepath.Join(root, "_posts", "post.ace"),
		filepath.Join(root, "_posts", "README"),
		filepath.Join(root, "other_dir", "post.md"),
	}
	if err := util.CreateEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the PostsCompiler,
	// the other files should be ignored.
	config.SourceDir = root
	config.PostsDir = filepath.Join(config.SourceDir, "_posts")
	PostsCompiler.Init()
	expectedPaths := []string{
		filepath.Join(root, "_posts", "post.md"),
	}

	// Use the MatchFunc to find all the paths
	config.SourceDir = root
	gotPaths, err := FindPaths(PostsCompiler.CompileMatchFunc())
	if err != nil {
		t.Error(err)
	}

	// Check that the paths we get are correct
	test_util.CheckStringsMatch(t, expectedPaths, gotPaths)
}

func TestPostsCompiler(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := string(os.PathSeparator) + filepath.Join("tmp", "test_posts_compiler")
	defer func() {
		// Remove everything after we're done
		if err := os.RemoveAll(root); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}
	}()

	// Copy some files from test_files to source directory in the temp root
	testFilesDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "albrow", "scribble", "test_files", "posts")
	srcDir := filepath.Join(root, "source")
	destDir := filepath.Join(root, "public")
	if err := util.RecursiveCopy(filepath.Join(testFilesDir, "source"), srcDir); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the posts
	config.SourceDir = filepath.Join(root, "source")
	config.PostsDir = filepath.Join(config.SourceDir, "_posts")
	config.LayoutsDir = filepath.Join(config.SourceDir, "_layouts")
	config.PostLayoutsDir = filepath.Join(config.SourceDir, "_post_layouts")
	config.DestDir = filepath.Join(root, "public")
	if err := PostsCompiler.Compile(filepath.Join(config.PostsDir, "post.md")); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	expectedDir := filepath.Join(testFilesDir, "public")
	expectedFile := filepath.Join(expectedDir, "post", "index.html")
	gotFile := filepath.Join(destDir, "post", "index.html")
	test_util.CheckFilesMatch(t, expectedFile, gotFile)
}
