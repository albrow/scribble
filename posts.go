package main

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/ace"
	"github.com/russross/blackfriday"
	"github.com/wsxiaoys/terminal/color"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

var (
	posts    = []*Post{}
	postsMap = map[string]*Post{} // a map of source path to post
)

type Post struct {
	Title       string        `toml:"title"`
	Author      string        `toml:"author"`
	Description string        `toml:"description"`
	Url         string        `toml:"-"` // the url for the post, not including protocol or domain name (useful for creating links)
	Content     template.HTML `toml:"-"` // the html content for the post (parsed from markdown source)
	dest        string        `toml:"-"` // the full destination path
	src         string        `toml:"-"` // the full source path
}

// parsePosts walks through postsDir, creates a new Post object for
// each markdown file there, and parses the content and frontmatter
// from the markdown files. It removes all the old posts and appends
// each new post created to posts and postsMap.
func parsePosts() {
	fmt.Printf("    reading posts in %s\n", postsDir)
	// remove any old posts
	posts = []*Post{}
	context["Posts"] = posts
	// walk through the source/posts dir
	if err := filepath.Walk(postsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// check if markdown file (ignore everything else)
		if filepath.Ext(path) == ".md" {
			// create a new Post object from the file and append it to posts
			p := createPostFromPath(path)
			p.parse()
		}
		return nil
	}); err != nil {
		panic(err)
	}
	fmt.Printf("    found %d posts\n", len(posts))
}

func createPostFromPath(path string) *Post {
	// create post object
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	p := &Post{
		Url:  "/" + name,
		dest: destDir + "/" + name,
		src:  path,
	}
	posts = append(posts, p)
	postsMap[path] = p
	context["Posts"] = posts
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
// for the post
func (p *Post) parse() {
	// open the source file
	file, err := os.Open(p.src)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)

	// split the file into frontmatter and markdown content
	frontMatter, content, err := split(r)
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

func getPostTemplate() *template.Template {
	// TODO: detect layout from frontmatter
	tpl, err := ace.Load("_layouts/base", "_views/post", &ace.Options{
		DynamicReload: true,
		BaseDir:       sourceDir,
		FuncMap:       funcMap,
	})
	if err != nil {
		chimeError(err)
	}
	return tpl
}

// compilePosts compiles all posts in posts
func compilePosts() {
	fmt.Println("--> compiling posts")
	tpl := getPostTemplate()
	for _, p := range posts {
		p.compileWithTemplate(tpl)
	}
}

// compile compiles a single post and writes it to
// the appropriate file in destDir
func (p Post) compile() {
	p.compileWithTemplate(getPostTemplate())
}

// compileWithTemplate compiles a single post using the given
// template and writes it to the appropriate file in destDir.
// It is useful for cases where you want to reuse the same
// template to compile more than one post in a for loop.
func (p Post) compileWithTemplate(tpl *template.Template) {
	// make the directory for the post
	err := os.Mkdir(p.dest, os.ModePerm)
	if err != nil {
		// if the directory already exists, that's fine
		// if there was some other error, panic
		if !os.IsExist(err) {
			panic(err)
		}
	}

	// make an index.html file inside that directory
	destPath := p.dest + "/index.html"
	color.Printf("@g    CREATE: %s -> %s\n", p.src, destPath)
	file, err := os.Create(destPath)
	if err != nil {
		// if the file already exists, that's fine
		// if there was some other error, panic
		if !os.IsExist(err) {
			panic(err)
		}
	}
	context["Post"] = p
	if err := tpl.Execute(file, context); err != nil {
		chimeError(err)
	}
	delete(context, "Post")
}

// remove removes all the files in destDir associated with the post
// and removes it from posts and postsMap
func (p *Post) remove() {
	if err := os.RemoveAll(p.dest); err != nil {
		if !os.IsNotExist(err) {
			// if the file doesn't exist that's fine,
			// otherwise throw an error
			panic(err)
		}
	}
	delete(postsMap, p.src)
	for i, other := range posts {
		if p == other {
			posts = append(posts[:i], posts[i+1:]...)
		}
	}
	fmt.Println("posts:", posts)
}
