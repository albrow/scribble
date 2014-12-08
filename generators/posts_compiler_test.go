package generators

import (
	"github.com/albrow/scribble/config"
	"os"
	"testing"
)

func TestPostsCompile(t *testing.T) {
	// Create a root path where all of our test files for this
	// test will live
	root := "/tmp/posts_compiler/"
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
		root + "_posts/post.md",
		root + "_posts/post.ace",
		root + "_posts/README",
		root + "other_dir/post.md",
	}
	if err := createEmptyFiles(tmpPaths); err != nil {
		t.Fatal(err)
	}

	// Only some paths are expected to be matched by the PostsCompiler,
	// the other files should be ignored.
	config.PostsDir = root + "_posts"
	PostsCompiler.Init()
	expectedPaths := []string{
		root + "_posts/post.md",
	}

	// Check that the paths we get are correct
	checkPathsMatch(t, PostsCompiler, root, expectedPaths)

	// Copy some files from test_files to the tmp directory
	gopath := os.Getenv("GOPATH")
	srcRoot := gopath + "/src/github.com/albrow/scribble/test_files/posts/"
	pathsToCopy := map[string]string{
		root + "_posts/post.md":    srcRoot + "_posts/post.md",
		root + "_views/post.ace":   srcRoot + "_views/post.ace",
		root + "_layouts/base.ace": srcRoot + "_layouts/base.ace",
	}
	if err := copyFiles(pathsToCopy); err != nil {
		t.Fatal(err)
	}

	// Attempt to compile the posts
	config.LayoutsDir = "_layouts"
	config.SourceDir = root
	config.PostsDir = "_posts"
	config.ViewsDir = "_views"
	if err := PostsCompiler.Compile(root+"_posts/post.md", root+"public"); err != nil {
		t.Fatal(err)
	}

	// Make sure the compiled result is correct
	checkOutputMatchesFile(t, root+"public/post/index.html", srcRoot+"expected.html")
}
