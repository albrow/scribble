Scribble
========

Version: 0.3.0

A tiny static blog generator written in go.

The core functionality for scribble is finished and it is pretty stable. However, there might
be some bugs and some things that will change before v1.0. I would caution anyone using scribble
for production-grade blogs, but I would encourage exploring the codebase and trying it out. v1.0
will be released in the near future and will be completely production-ready.

If you encounter any issues or have any suggestions, please open a pull request :)


How it Works
------------

Scribble is a command line tool that compiles different types of source files into a
static blog made up of html, css, and js. It uses:

- Markdown for writing posts
- Toml for frontmatter and configuration
- Sass for styling
- Standard go html templates for pages

Scribble is optimized for speed and usability. It compiles the source files for a medium-sized blog
in next to no time. It uses sassc (a C port of the sass compiler) to compile sass. It also features
a built in server and can automatically recompile whenever you change files.


Quickstart Guide
----------------

### Basic Commands

You can run `scribble --help` to see a description of commands and flags. You can also run
`scribble help <command>` to see a more detailed description of a specific command, including
the flags it supports. The other commands are:

- `version`: print the version number.
- `compile`: compile your blog into static html, css, and javascript. The `-w` flag will tell
	scribble to watch for changes and recompile automatically.
- `serve`: compile and serve your blog; also watches for changes and recompiles automatically.


### File Structure

I've created a [seed project](https://github.com/albrow/scribble-seed) for scribble which you can
clone to get a basic blog working. (In the future, there may be a `scribble new` command which will
automatically create this skeleton for you). To clone the seed project, just run
`git clone https://github.com/albrow/scribble-seed`.

The file structure of the seed project looks like this:

```
blog
├── config.toml
├── public
└── source
    ├── _includes
    │   ├── foot.tmpl
    │   └── head.tmpl
    ├── _layouts
    │   └── base.tmpl
    ├── _post_layouts
    │   └── post.tmpl
    ├── _posts
    │   ├── one.md
    │   ├── three.md
    │   └── two.md
    ├── index.tmpl
    ├── js
    │   └── main.js
    └── styles
        ├── _colors.scss
        ├── _fonts.scss
        └── main.scss
```

- `config.toml` is required and stores some basic configuration in [toml](https://github.com/toml-lang/toml)
format. This consists of metadata such as your Blog's title, the author's name, and a description, but also
tells scribble where to look for certain files. Every scribble project must have a config.toml file
in the project root directory.

- `public` is the folder where scribble will put your finished website after compiling. It's also the
folder that scribble will serve from when using the `scribble serve` command. This is set via the
`destDir` key in `config.toml`. `destDir` is required but you can set it to anything you want.

- `source` is where all the source code you write will live. This includes things like stylesheets,
posts, html templates, and javascript. This is set via the `sourceDir` key in `config.toml`.
`sourceDir` is required but you can set it to anything you want.

- `source/_includes` is an optional folder where you can put partial html templates, i.e. templates
which don't constitute a full page on their own, but are meant to be *included* in other templates.
In the seed project, there are two files in `source/_includes`, one for filling in the `<head>`
tag with metadata and stylesheets, and one for including any javascript files at the bottom of the
`<body>` tag. This is set via the `includesDir` key in `config.toml`. `includesDir` is optional and
you don't have to have any includes if you don't want them.

- `source/_layouts` is where you put html layouts. Layouts are reusable html wrappers that define how
certain pages will look. In the seed project, there is just one layout, called `base.tmpl`. It consists
of html boilerplate like the `<html>`, `<head>`, and `<body>` tags. It includes the two files in our
`_includes` directory by using the html/template syntax `{{ template "head.tmpl" . }}`. The layouts
directory is defined via the `layoutsDir` key in `config.toml`. `layoutsDir` is required, along with at
least one layout but you can set it to anything you want.

- `source/_post_layouts` is where you put post layouts. Like html layouts, post layouts are reusable
html wrappers that define how certain posts will look. In the seed project, there is just one post layout,
called `post.tmpl`. The default template simply consists of the title of the post in a header and
the content of the post wrapped in a div below it. The post layouts directory is defined via the
`postLayoutsDir` key in `config.toml`. `postLayoutsDir` is required, along with at least one post layout,
but you can set it to anything you want.

- `_posts` is where your posts will reside. Posts are written in markdown and must include toml
frontmatter which defines the post layout, and optionally the title, author's name, date, and
description. The posts directory is defined via the `postsDir` key in `config.toml`. `postsDir` is
required but you can set it to anything you want.

- `index.tmpl` is the index page and will compile to index.html. It consists of an unordered list
of links to the 5 most recent posts. It also has frontmatter at the top of the file to tell scribble
to use the `base.tmpl` layout. You are not required to have an `index.tmpl` file, and you can organize
your pages however you want (see compilation details below). 

- `js` is a folder where you can put javascript. Any files here will be copied to `destDir` directly.
You can actually put your javascript files in a different folder if you want. Scribble will
pick them up no matter what directory they're in (see compilation details below). If your blog doesn't
need any javascript, you can simply remove this folder.

- `styles` is where you can put sass stylesheets. They will be picked up by the sass compiler, and
any sass files that don't start with an underscore will be compiled to css files in `destDir`. Just
like javascript files, you can just put sass files in a different directory if you want. Scribble will
pick them up no matter what directory they're in (see compilation details below).


### Compilation

In general, files and folders that start with an underscore have special meaning and will be ignored
(i.e. not copied over to `destDir`). You can use this fact to prevent things like partial templates or
sass imports from being published. More specifically, compilation follows these rules:

1. Any markdown files (identified by the .md extension) in `postsDir` get treated as posts and are
	converted to html. Specifically, They are converted to an index.html file in a folder with the
	same name as the markdown file. So `source/_posts/first.md` becomes `public/first/index.html`.
	They are also added to an in-memory representation of posts and their metadata is accessible through
	the [Posts function](https://github.com/albrow/scribble/blob/dc25cd04f111659d19cd8b9456488a949a79aedd/compilers/posts_compiler.go#L193)
	in templates. Markdown files anywhere else are currently ignored, but may be converted to html in future
	versions. That's why your `postsDir` should start with an underscore, so that your posts will be distinct
	from markdown pages.
2. Any sass files (identified by the .scss extension) that do not start with an underscore and are 
	not in a directory that starts with an underscore are converted to css and retain the same filename
	and relative path. So `source/styles/base/main.scss` becomes `public/styles/base/main.css`, and
	`source/styles/base/_fonts.scss` is not copied over to `destDir`.
3. Any go html template files (identified by the .tmpl extension) that do not start with an underscore
	and are not in a directory that starts with an underscore are converted to html and retain the same
	filename and relative path. So `source/about/index.tmpl` becomes `public/about/index.html`, and
	`source/about/_partials.tmpl` is not copied over to `destDir`. This is why your `layoutsDir`,
	`includesDir`, and `postLayoutsDir` should have names that start with underscores, because we don't
	want those files to be directly converted to html.
4. Any other files in `sourceDir` that do not start with an underscore and are not in a directory that
	starts with an underscore are simply copied over as is, retaining their relative paths. So
	`source/js/main.js` becomes `public/js/main.js` and `source/js/_libs/watch.js` would not
	be copied over to `destDir`.

If you compile the default seed project, the compiled blog would look like this:

```
public
├── index.html
├── js
│   └── main.js
├── one
│   └── index.html
├── styles
│   └── main.css
├── three
│   └── index.html
└── two
    └── index.html
```

- `index.html` came from `source/index.tmpl`

- `js/main.js` came from `source/js/main.js`

- `one/index.html` came from `source/_posts/one.md`, and the other posts - `two` and `three` came from their
respective source files.

- `styles/main.css` came from `source/styles/main.scss`

Note that the `_includes`, `_layouts`, and `_post_layouts` folders were not compiled because they started
with underscores. Same with the imported sass files: `_fonts.scss` and `_colors.scss`.

### Posts

Posts are markdown files located in `postsDir`. Post files should have toml frontmatter defining the layout
that should be used, as well as other metadata such as the title, date, author, and description. In the future,
you will be able to add your own metadata, but for now this is all that is supported.

Here's an example of a simple post file with all the frontmatter included:

``` markdown
+++
title = "My First Blog Post"
author = "Your Name"
date = "2014-11-16T13:50:53-05:00"
layout = "post.tmpl"
description = "The first post I've ever written using scribble!"
+++

This is a post.

### This is a header

This is a paragraph.
```

#### Related Resources:

- [Learn more about toml](https://github.com/toml-lang/toml).
- [Learn more about markdown](http://daringfireball.net/projects/markdown/).

### Sass

Any sass files that have the .scss extension will be compiled into css automatically (unless they start
with an underscore as described above in the Compilation section). If you already know sass, you don't have
to change anything about the way you write sass files.

#### Related Resources:

[Learn more about sass](http://sass-lang.com/).

### Html Templates

Scribble uses go's [html/template](http://golang.org/pkg/html/template/) package for html templates. This
includes regular html pages, as well as layouts and includes. Admittedly, this may be the component that is
the hardest to grok at first. I strongly recommend the related resources at the bottom of this section.

All templates should have the .tmpl extension. Scribble uses the following conventions, but you are not strictly
required to adhere to them:

1. Html layouts (which exist in `layoutsDir`) should be a wrapper around a named template, typically called
`"content"`. Templates which use the layout will be required to define the content template.
2. Layouts and includes are identified by their base filename, including the extension. (This is the go
default when using `ParseFiles` or `ParseGlob`).
3. Post layouts define a content template and may use a layout as a wrapper.

You are however, required to define a layout in the frontmatter of any html template. The layout file is
associated with the html template when parsing it.

Here are some examples taken from the [seed project](https://github.com/albrow/scribble-seed).

`_layouts/base.tmpl`: The default layout.
``` html
<!DOCTYPE html>
<html>
<head>
	<!-- 
		"head.tmpl" references a file in the _includes folder.
		The . following the template name passes the current context through.
	-->
	{{ template "head.tmpl" . }}
</head>
<body>
	<!-- 
		"content" references a named template which html templates should define
		using the define keyword. This is similar to "yeild" in other templating
		languages.
	-->
	{{ template "content" . }}
	<!-- 
		"foot.tmpl" is a file in the _includes folder, just like "head.tmpl".
	-->
	{{ template "foot.tmpl" . }}
</body>
</html>
```

`index.tmpl`: The index page (gets converted to index.html).
``` html
+++
# layout is required to be defined in the frontmatter
layout = "base.tmpl"
+++
<!-- 
	We wrap the entire template in a define "content" block so that when
	the templates are executed and compiled, everything here will take the
	place of the corresponding {{ template "content" . }} block in base.tmpl
-->
{{ define "content" }}
<h1>Recent Posts</h1>
<ul>
	<!-- 
		Posts is a function accessible in any template. It returns a number of
		posts sorted by date. It takes one argument, the maximum number of posts
		to return.
	-->
	{{range Posts 5}}
		<li>
			<!-- 
				Render a link to each post.
			-->
			<a href="{{.Url}}">{{.Title}}</a>
		</li>
	{{end}}
</ul>
{{ end }}
```

`_post_layouts/post.tmpl`: The default layout for posts.
``` html
<!-- 
	We wrap the entire template in a define "content" block, just like we did
	in index.tmpl.
-->
{{ define "content" }}
<a href="/">&larr; Back to Home</a>
<!-- 
	Simply render the title of the post in a header and the content wrapped
	in a div.
-->
<h2 class="post-title">{{.Post.Title}}</h2>
<div class="post-content">
	{{.Post.Content}}
</div>
{{ end }}
<!-- 
	Execute base.tmpl, where everything inside the content template will be rendered.
-->
{{ template "base.tmpl" . }}
```

#### Related Resources:

- [Learn more about go's text templates](http://golang.org/pkg/text/template/), which share a lot of functionality with html templates.
- [Learn more about go's html templates](http://golang.org/pkg/html/template/).
- [Learn more about layout/template inheritance in go](https://elithrar.github.io/article/approximating-html-template-inheritance/). (Scribble takes inspiration from the ideas there).
