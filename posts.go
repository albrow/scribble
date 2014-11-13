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

var posts = []Post{}

type Post struct {
	Title       string        `toml:"title"`
	Author      string        `toml:"author"`
	Description string        `toml:"description"`
	Content     template.HTML `toml:"-"`
	Url         string        `toml:"-"`
	Dir         string        `toml:"-"`
}

func parsePosts() {
	fmt.Printf("    parsing posts in %s\n", postsDir)
	// remove any old posts
	posts = []Post{}
	context["Posts"] = posts
	// walk through the source/posts dir
	if err := filepath.Walk(postsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// check if markdown file (ignore everything else)
		if filepath.Ext(path) == ".md" {
			// create a new Post object from the file and append it to posts
			p, err := createPostFromPath(path, info)
			if err != nil {
				return err
			}
			posts = append(posts, p)
		}
		context["Posts"] = posts
		return nil
	}); err != nil {
		panic(err)
	}
	fmt.Printf("    found %d posts\n", len(posts))
}

func createPostFromPath(path string, info os.FileInfo) (Post, error) {
	// create post object
	name := strings.TrimSuffix(info.Name(), filepath.Ext(path))
	p := Post{
		Url: "/" + name,
		Dir: name,
	}

	// open the source file
	file, err := os.Open(path)
	if err != nil {
		return p, err
	}

	// extract and parse front matter
	if err := p.parseFromFile(file); err != nil {
		return p, err
	}

	return p, nil
}

func (p *Post) parseFromFile(file *os.File) error {
	r := bufio.NewReader(file)
	frontMatter, content, err := split(r)
	if err != nil {
		return err
	}
	if _, err := toml.Decode(frontMatter, p); err != nil {
		return err
	}
	p.Content = template.HTML(blackfriday.MarkdownCommon([]byte(content)))
	return nil
}

func compilePosts() {
	fmt.Println("    compiling posts")
	// TODO: detect layout from frontmatter
	// load the template
	tpl, err := ace.Load("_layouts/base", "_views/post", &ace.Options{
		DynamicReload: true,
		BaseDir:       sourceDir,
		FuncMap:       funcMap,
	})
	if err != nil {
		chimeError(err)
	}
	for _, p := range posts {
		p.compile(tpl)
	}
}

func (p Post) compile(tpl *template.Template) {
	dirName := destDir + "/" + p.Dir

	// make the directory for each post
	err := os.Mkdir(dirName, os.ModePerm|os.ModeDir)
	if err != nil {
		// if the directory already exists, that's fine
		// if there was some other error, panic
		if !os.IsExist(err) {
			panic(err)
		}
	}

	// make an index.html file inside that directory
	destPath := dirName + "/index.html"
	color.Printf("@g    CREATE: %s\n", destPath)
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
