package compilers

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/context"
	"github.com/albrow/scribble/util"
	"github.com/howeyc/fsnotify"
	"github.com/russross/blackfriday"
	"github.com/wsxiaoys/terminal/color"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PostsCompilerType represents a type capable of compiling post files.
type PostsCompilerType struct {
	pathMatch string
	// createdDirs keeps track of the directories that were created in config.DestDir.
	// It is used in the RemoveOld method.
	createdDirs []string
}

// PostCompiler is an instatiation of PostCompilerType
var PostsCompiler = PostsCompilerType{
	pathMatch: "",
}

// Post is an in-memory representation of the metadata for a given post.
// Much of this data comes from the toml frontmatter.
type Post struct {
	Title       string    `toml:"title"`
	Author      string    `toml:"author"`
	Description string    `toml:"description"`
	Date        time.Time `toml:"date"`
	// the url for the post, not including protocol or domain name (useful for creating links)
	Url template.URL `toml:"-"`
	// the html content for the post (parsed from markdown source)
	Content template.HTML `toml:"-"`
	// the full source path
	src string `toml:"-"`
	// the layout tmpl file to be used for the post
	Layout string `toml:"layout"`
	// the html template that can be used to render the post
	template *template.Template
}

var (
	// a slice of all posts
	posts = []*Post{}
	// a map of source path to post
	postsMap = map[string]*Post{}
)

// Init should be called before any other methods. In this case, Init
// sets up the pathMatch variable based on config.SourceDir and config.PostsDir
// and adds the Posts helper function to FuncMap.
func (p *PostsCompilerType) Init() {
	p.pathMatch = filepath.Join(config.PostsDir, "*.md")
	// Add the posts function to FuncMap
	context.FuncMap["Posts"] = Posts
}

// CompileMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is any file that is inside config.PostsDir and ends in ".md", excluding
// hidden files and directories (which start with a ".") but not those
// which start with an underscore.
func (p *PostsCompilerType) CompileMatchFunc() MatchFunc {
	return pathMatchFunc(p.pathMatch, true, false)
}

// WatchMatchFunc returns a MatchFunc which will return true for
// any files which match a given pattern. In this case, the pattern
// is the same as it is for CompileMatchFunc.
func (p *PostsCompilerType) WatchMatchFunc() MatchFunc {
	// PostsCompiler needs to watch all posts in the posts dir,
	// but also needs to watch all the *tmpl files in the posts
	// layouts dir, because if those change, it affects the way
	// posts are rendered.
	postsMatch := pathMatchFunc(p.pathMatch, true, false)
	postLayoutsMatch := pathMatchFunc(filepath.Join(config.PostLayoutsDir, "*.tmpl"), true, false)
	// unionMatchFuncs combines these two cases and returns a MatchFunc
	// which will return true if either matches. This allows us to watch
	// for changes in both the posts dir and the posts layouts dir.
	return unionMatchFuncs(postsMatch, postLayoutsMatch)
}

// Compile compiles the file at srcPath. The caller will only
// call this function for files which belong to PostsCompiler
// according to the MatchFunc. Behavior for any other file is
// undefined. Compile will output the compiled result to the appropriate
// location in config.DestDir.
func (p *PostsCompilerType) Compile(srcPath string) error {
	// Get the parsed post object and determine dest path
	post := getOrCreatePostFromPath(srcPath)
	srcFilename := filepath.Base(srcPath)
	destPath := fmt.Sprintf("%s/%s", config.DestDir, strings.TrimSuffix(srcFilename, ".md"))
	destIndexFilePath := filepath.Join(destPath, "index.html")
	color.Printf("@g    CREATE: %s -> %s\n", srcPath, destIndexFilePath)

	// Create the index file
	destFile, err := util.CreateFileWithPath(destIndexFilePath)
	if err != nil {
		return err
	}

	// Parse content and frontmatter, then set the appropriate layout based on
	// the layout key in the frontmatter
	if err := post.parse(); err != nil {
		return err
	}

	// Write to destFile by executing the template
	postContext := context.CopyContext()
	postContext["Post"] = post
	if err := post.template.Execute(destFile, postContext); err != nil {
		return fmt.Errorf("ERROR compiling html template for posts: %s", err.Error())
	}

	// Add the created dir to the list of created dirs
	p.createdDirs = append(p.createdDirs, destPath)

	return nil
}

// CompileAll compiles zero or more files identified by srcPaths.
// It works simply by calling Compile for each path. The caller is
// responsible for only passing in files that belong to AceCompiler
// according to the MatchFunc. Behavior for any other file is undefined.
func (p *PostsCompilerType) CompileAll(srcPaths []string) error {
	fmt.Println("--> compiling posts")
	for _, srcPath := range srcPaths {
		if err := p.Compile(srcPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *PostsCompilerType) FileChanged(srcPath string, ev fsnotify.FileEvent) error {
	// Because of the way we set up the watcher, there are two possible
	// cases here.
	// 1) A template in the post layouts dir was changed. In this case,
	// we would ideally recompile all the posts that used that layout.
	// 2) A markdown file corresponding to a single post was changed. In this
	// case, ideally we only recompile the post that was changed. We need to
	// take into account any rename, delete, or create events and how they
	// affect the output files in destDir.
	switch filepath.Ext(srcPath) {
	case ".md":
		// TODO: Be more intelligent here? If a single post file was midified,
		// we can simply recompile that post. We would also need to take into
		// account the subtle differences between rename, create, and delete
		// events. For now, recompile all posts.
		if err := recompileAllForCompiler(p); err != nil {
			return err
		}
		return nil
	case ".tmpl":
		// TODO: Analyze post files and be more intelligent here?
		// When a post layout changes, only recompile the posts that
		// use that layout. For now, recompile all posts.
		if err := recompileAllForCompiler(p); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (p *PostsCompilerType) RemoveOld() error {
	// Simply iterate through createdDirs and remove each of them
	// NOTE: this is different from the other compilers because each
	// post gets created as an index.html file inside some directory
	// (for prettier urls). So instead of removing files, we're removing
	// directories.
	for _, dir := range p.createdDirs {
		if err := util.RemoveAllIfExists(dir); err != nil {
			return err
		}
	}
	return nil
}

// Posts returns up to limit posts. If limit is 0, it returns
// all posts. If limit is greater than len(posts), it returns
// all posts.
func Posts(limit ...int) []*Post {
	if len(limit) == 0 || limit[0] == 0 || limit[0] > len(posts) {
		return posts
	} else {
		return posts[:limit[0]]
	}
}

func createPostFromPath(path string) *Post {
	// create post object
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	p := &Post{
		Url: template.URL("/" + name),
		src: path,
	}
	posts = append(posts, p)
	postsMap[path] = p
	return p
}

func getPostByPath(path string) *Post {
	return postsMap[path]
}

func getOrCreatePostFromPath(path string) *Post {
	if p, found := postsMap[path]; found {
		return p
	} else {
		return createPostFromPath(path)
	}
}

// parse reads from the source file and sets the content and metadata fields
// for the post. It also creates a new template for the post using the layout
// field of the frontmatter.
func (p *Post) parse() error {
	// Open the source file
	file, err := os.Open(p.src)
	if err != nil {
		return err
	}
	r := bufio.NewReader(file)

	// Split the file into frontmatter and markdown content
	frontMatter, content, err := util.SplitFrontMatter(r)
	if err != nil {
		return err
	}

	// Decode the frontmatter
	if _, err := toml.Decode(frontMatter, p); err != nil {
		return err
	}

	// Parse the markdown content and set p.Content
	p.Content = template.HTML(blackfriday.MarkdownCommon([]byte(content)))

	// Create a template for the post
	if p.Layout == "" {
		return fmt.Errorf("Could not find layout definition in toml frontmatter for post: %s", p.src)
	}

	// When we call template.ParseFiles, we want the post layout file to be first,
	// so it is the one that will be executed.
	postLayoutFile := filepath.Join(config.PostLayoutsDir, p.Layout)
	otherLayoutFiles, err := filepath.Glob(filepath.Join(config.LayoutsDir, "*.tmpl"))
	if err != nil {
		return err
	}
	allFiles := append([]string{postLayoutFile}, otherLayoutFiles...)
	// Parse the includes files if there are any
	if config.IncludesDir != "" {
		includeFiles, err := filepath.Glob(filepath.Join(config.IncludesDir, "*tmpl"))
		if err != nil {
			return err
		}
		allFiles = append(allFiles, includeFiles...)
	}
	tmpl, err := template.ParseFiles(allFiles...)
	if err != nil {
		return err
	}
	p.template = tmpl

	return nil
}

// The PostsByDate type is used only for sorting
type PostsByDate []*Post

func (p PostsByDate) Len() int {
	return len(p)
}

func (p PostsByDate) Less(i, j int) bool {
	return p[i].Date.Before(p[j].Date)
}

func (p PostsByDate) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
