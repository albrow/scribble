package lib

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/russross/blackfriday"
	"github.com/yosssi/ace"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	posts    = []*Post{}
	postsMap = map[string]*Post{} // a map of source path to post
)

type Post struct {
	Title       string        `toml:"title"`
	Author      string        `toml:"author"`
	Description string        `toml:"description"`
	Date        time.Time     `toml:"date"`
	Url         string        `toml:"-"` // the url for the post, not including protocol or domain name (useful for creating links)
	Content     template.HTML `toml:"-"` // the html content for the post (parsed from markdown source)
	src         string        `toml:"-"` // the full source path
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

// ParsePosts walks through PostsDir, creates a new Post object for
// each markdown file there, and parses the content and frontmatter
// from the markdown files. It removes all the old posts and appends
// each new post created to posts and postsMap.
func ParsePosts() {
	fmt.Printf("    reading posts in %s\n", PostsDir)
	// remove any old posts
	posts = []*Post{}
	// walk through the source/posts dir
	if err := filepath.Walk(PostsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// check if markdown file (ignore everything else)
		if filepath.Ext(path) == ".md" {
			// create a new Post object from the file and append it to posts
			p := CreatePostFromPath(path)
			p.Parse()
		}
		return nil
	}); err != nil {
		panic(err)
	}
	// Sort the posts
	sort.Stable(PostsByDate(posts))
	fmt.Printf("    found %d posts\n", len(posts))
}

func CreatePostFromPath(path string) *Post {
	// create post object
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	p := &Post{
		Url: "/" + name,
		src: path,
	}
	posts = append(posts, p)
	postsMap[path] = p
	return p
}

func GetPostByPath(path string) *Post {
	return postsMap[path]
}

func GetOrCreatePostFromPath(path string) *Post {
	if p, found := postsMap[path]; found {
		return p
	} else {
		return CreatePostFromPath(path)
	}
}

// Parse reads from the source file and sets the content and metadata fields
// for the post
func (p *Post) Parse() {
	// open the source file
	file, err := os.Open(p.src)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)

	// split the file into frontmatter and markdown content
	frontMatter, content, err := SplitFrontMatter(r)
	if err != nil {
		panic(err)
	}

	// decode the frontmatter
	if _, err := toml.Decode(frontMatter, p); err != nil {
		panic(err)
	}

	// parse the markdown content and set p.Content
	p.Content = template.HTML(blackfriday.MarkdownCommon([]byte(content)))
}

// GetPostTemplate returns the ace template which
// is to be used for rendering all posts.
func GetPostTemplate() *template.Template {
	// TODO: detect layout from frontmatter
	tpl, err := ace.Load("_layouts/base", ViewsDir+"/post", &ace.Options{
		DynamicReload: true,
		BaseDir:       SourceDir,
		FuncMap:       FuncMap,
	})
	if err != nil {
		ChimeError(err)
	}
	return tpl
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
