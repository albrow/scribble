package generators

import (
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/util"
	"os"
	"testing"
)

func TestPostsPathMatch(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/posts_compiler_paths"
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
		root + "/_posts/post.md",
		root + "/_posts/post.ace",
		root + "/_posts/README",
		root + "/other_dir/post.md",
	}
	if err := createEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the PostsCompiler,
	// the other files should be ignored.
	config.SourceDir = root
	config.PostsDir = "_posts"
	PostsCompiler.Init()
	expectedPaths := []string{
		root + "/_posts/post.md",
	}

	// Check that the paths we get are correct
	checkPathsMatch(t, PostsCompiler, root, expectedPaths)
}

func testPostsCompiler(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/test_posts_compiler"
	defer func() {
		// Remove everything after we're done
		if err := os.RemoveAll(root); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}
	}()

	// Copy some files from test_files to source directory in the temp root
	testFilesDir := os.Getenv("GOPATH") + "/src/github.com/albrow/scribble/test_files/posts"
	srcDir := root + "/source"
	destDir := root + "/public"
	if err := util.RecursiveCopy(testFilesDir+"/source", srcDir); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the posts
	config.SourceDir = root
	config.PostsDir = "_posts"
	config.LayoutsDir = "_layouts"
	config.ViewsDir = "_views"
	config.DestDir = root + "/public"
	if err := PostsCompiler.Compile(root + "_posts/post.md"); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	expectedDir := testFilesDir + "/public"
	checkOutputMatchesFile(t, destDir+"/post/index.html", expectedDir+"/post/index.html")
}
