package generators

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/ace"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/context"
	"github.com/albrow/scribble/util"
	"github.com/russross/blackfriday"
	"github.com/wsxiaoys/terminal/color"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PostsCompilerType struct {
	pathMatch string
}

var PostsCompiler = PostsCompilerType{
	pathMatch: "",
}

type Post struct {
	Title       string        `toml:"title"`
	Author      string        `toml:"author"`
	Description string        `toml:"description"`
	Date        time.Time     `toml:"date"`
	Url         string        `toml:"-"` // the url for the post, not including protocol or domain name (useful for creating links)
	Content     template.HTML `toml:"-"` // the html content for the post (parsed from markdown source)
	src         string        `toml:"-"` // the full source path
}

var (
	posts    = []*Post{}
	postsMap = map[string]*Post{} // a map of source path to post
)

func (p *PostsCompilerType) Init() {
	p.pathMatch = config.PostsDir + "/*.md"
}

func (p PostsCompilerType) GetWalkFunc(paths *[]string) filepath.WalkFunc {
	return pathMatchWalkFunc(paths, p.pathMatch, true, false)
}

func (p PostsCompilerType) Compile(srcPath string, destDir string) error {
	// Get the parsed post object and determine dest path
	post := getOrCreatePostFromPath(srcPath)
	srcFilename := filepath.Base(srcPath)
	destPath := fmt.Sprintf("%s/%s", destDir, strings.TrimSuffix(srcFilename, ".md"))
	destIndexFilePath := destPath + "/index.html"
	color.Printf("@g    CREATE: %s -> %s\n", srcPath, destIndexFilePath)

	// Create the index file
	fmt.Println("creating file")
	destFile, err := util.CreateFileWithPath(destIndexFilePath)
	if err != nil {
		return err
	}

	// Get and compile the template
	fmt.Println("getting template")
	tmpl := getPostTemplate()
	fmt.Println("parsing post")
	post.parse()
	fmt.Println("getting context")
	postContext := context.GetContext()
	postContext["Post"] = post
	fmt.Println("Executing template")
	if err := tmpl.Execute(destFile, postContext); err != nil {
		return fmt.Errorf("ERROR compiling ace template for posts: %s", err.Error())
	}
	fmt.Println("Done")
	return nil
}

func (p PostsCompilerType) CompileAll(srcPaths []string, destDir string) error {
	for _, srcPath := range srcPaths {
		if err := p.Compile(srcPath, destDir); err != nil {
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
		Url: "/" + name,
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
// for the post
func (p *Post) parse() {
	// open the source file
	file, err := os.Open(p.src)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)

	// split the file into frontmatter and markdown content
	frontMatter, content, err := util.SplitFrontMatter(r)
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

// getPostTemplate returns the ace template which
// is to be used for rendering all posts.
func getPostTemplate() *template.Template {
	// TODO: detect layout from frontmatter
	tpl, err := ace.Load("_layouts/base", config.ViewsDir+"/post", &ace.Options{
		DynamicReload: true,
		BaseDir:       config.SourceDir,
		FuncMap:       context.FuncMap,
	})
	if err != nil {
		util.ChimeError(err)
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
