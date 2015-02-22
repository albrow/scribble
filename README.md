Scribble
========

Version: 0.4.0

A tiny static blog generator written in go.

The core functionality for scribble is finished and it is pretty stable. However, there might
be some bugs and some things that will change before v1.0. I would caution anyone using scribble
for production-grade blogs, but I would encourage exploring the codebase and trying it out. v1.0
will be released in the near future and will be completely production-ready.

If you encounter any issues or have any suggestions, please open a pull request :)


Installation
------------

### Prerequisites

If you want scribble to compile sass for you, you  must install sassc, which is a C port
of the sass library. You may need to [download and install sassc from source](https://github.com/sass/sassc).
Future versions of scribble may relax this requirement and fallback on the ruby implementation
if you have it. If you don't want to use sass, you don't have to install it.

As of scribble v0.4.0, you can use either
[go's native html templates](http://golang.org/pkg/html/template/) or [jade](http://jade-lang.com/)
for pages and layouts. If you want to use jade, you need to
[install it first](https://github.com/jadejs/jade#installation). Jade requires node, so if
you don't have node you will need to [install that as well](http://nodejs.org/). After you
are done installing, make sure the jade executable is in your PATH and that you can run `jade`
from the command line. If you plan to use go's native html templates, you don't have to install
anything. 

### Pkg Installer

If you are running mac os x 10.5+, you can download a pkg file from the
[Releases page](https://github.com/albrow/scribble/releases) which will guide you through
the process of installing scribble automatically. For other platforms, you will need to
install via go get for now.

### Install via Go Get

1. [Download and install go](https://golang.org/dl/).
2. Follow [these instructions](https://golang.org/doc/code.html) for setting up your go workspace.
3. Run `go get -u github.com/albrow/scribble`. To clone the latest version of scribble and install it into `$GOPATH/bin`.
4. If you have added `$GOPATH/bin` to your `$PATH`, you can run scribble directly. Try running scribble with `scribble version`.


How it Works
------------

Scribble is a command line tool that compiles different types of source files into a
static blog made up of html, css, and js. It uses:

- Markdown for writing posts
- Toml for frontmatter and configuration
- Sass for styling
- Standard go html templates or jade for pages and layouts

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
clone to get a basic blog working. In the future, there may be a `scribble new` command which will
automatically create this skeleton for you, but for now you will need to clone the repository. On
a unix-like system, you can run the following command to clone the version of scribble-seed
corresponding to your scribble version:

``` bash
git clone -b `scribble version` --depth 1 https://github.com/albrow/scribble-seed.git
```

As of v0.4.0, scribble-seed uses jade as the default templating language. If you want to use
go's native html templates instead, see the section on [Html Templates](#html-templates) below.

The file structure of the seed project looks like this:

```
blog
├── config.toml
├── public
└── source
    ├── _includes
    │   ├── foot.jade
    │   └── head.jade
    ├── _layouts
    │   └── base.jade
    ├── _post_layouts
    │   └── post.jade
    ├── _posts
    │   ├── one.md
    │   ├── three.md
    │   └── two.md
    ├── index.jade
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

- `source/_includes` is an optional folder where you can put partial templates, i.e. templates
which don't constitute a full page on their own, but are meant to be *included* in other templates.
In the seed project, there are two files in `source/_includes`, one for filling in the `<head>`
tag with metadata and stylesheets, and one for including any javascript files at the bottom of the
`<body>` tag. If you are using go's native html/templates, you must tell scribble where the includes
are located via the `includesDir` key in `config.toml`. If you are using jade, the `includesDir` key
has no effect and you must specifiy includes by their full or relative paths. However it's still a good
idea to organize your includes into a single directory.

- `source/_layouts` is where you put html layouts. Layouts are reusable wrappers that define how
certain pages will look. In the seed project, there is just one layout, called `base.jade`. It consists
of html boilerplate like the `<html>`, `<head>`, and `<body>` tags. It includes the two files in our
`_includes` directory. If you are using go's native templates, the layouts directory is required and 
must be defined via the `layoutsDir` key in `config.toml`. However, if you are using jade, the `layoutsDir`
key has no effect and you must reference layouts by their full or relative paths. However it's still a
good idea to organize your layouts into a single directory.

- `source/_post_layouts` is where you put post layouts. Like html layouts, post layouts are reusable
html wrappers that define how certain posts will look. In the seed project, there is just one post layout,
called `post.jade`. The default template simply consists of the title of the post in a header and
the content of the post wrapped in a div below it. The post layouts directory is defined via the
`postLayoutsDir` key in `config.toml`. `postLayoutsDir` is required, along with at least one post layout.
Post layouts can be a `.tmpl` file if you wish to use go's native templates instead of jade.

- `_posts` is where your posts will reside. Posts are written in markdown and must include toml
frontmatter which defines the post layout, and optionally the title, author's name, date, and
description. The posts directory is defined via the `postsDir` key in `config.toml`. `postsDir` is
required but you can set it to anything you want.

- `index.jade` is the index page and will compile to index.html. It consists of an unordered list
of links to the 5 most recent posts. You are not required to have an `index.jade` file, and you can
organize your pages however you want (see compilation details below). 

- `js` is a folder where you can put javascript. Any files here will be copied to `destDir` directly.
You can actually put your javascript files in a different folder if you want. Scribble will
pick them up no matter what directory they're in (see compilation details below). If your blog doesn't
need any javascript, you can simply remove this folder.

- `styles` is where you can put sass stylesheets. They will be picked up by the sass compiler, and
any sass files that don't start with an underscore will be compiled to css files in `destDir`. Just
like javascript files, you can just put sass files in a different directory if you want. Scribble will
pick them up no matter what directory they're in (see compilation details below). If you don't want
to use sass at all, you can put css files here and they will also be picked up by the compiler.


### Compilation

In general, files and folders that start with an underscore have special meaning and will be ignored
(i.e. not copied over to `destDir`). You can use this fact to prevent things like partial templates or
sass imports from being published. More specifically, compilation follows these rules:

1. Any markdown files (identified by the .md extension) in `postsDir` get treated as posts and are
	converted to html. Specifically, They are converted to an index.html file in a folder with the
	same name as the markdown file. So `source/_posts/first.md` becomes `public/first/index.html` and
	can be accessed by the url `public/first`. They are also added to an in-memory representation of
	posts and their metadata is accessible through the [Posts function](https://github.com/albrow/scribble/blob/dc25cd04f111659d19cd8b9456488a949a79aedd/compilers/posts_compiler.go#L193) if you are using go's
	native templates, or the `Posts` key if you are using jade. Markdown files anywhere else are
	currently ignored, but may be converted to html in future versions. That's why your `postsDir`
	should start with an underscore, so that your posts will be distinct from markdown pages.
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
4. Any jade files (identified by the .jade extension) that do not start with an underscore and are
	not in a directory that starts with an underscore are converted to html and retain the same
	filename and relative path. So `source/about/index.jade` becomes `public/about/index.html`, and
	`source/about/_partials.jade` is not copied over to `destDir`. This is why your `layoutsDir`,
	`includesDir`, and `postLayoutsDir` should have names that start with underscores, because we don't
	want those files to be directly converted to html.
5. Any other files in `sourceDir` that do not start with an underscore and are not in a directory that
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

- `index.html` came from `source/index.jade`

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
layout = "post.jade"
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
to change anything about the way you write sass with scribble.

#### Related Resources:

[Learn more about sass](http://sass-lang.com/).

### Jade

You may use the [jade templating language](http://jade-lang.com/) for html templates and layouts.
If you already know jade, you don't have to change anything about the way that you write jade with
scribble.

#### Related Resources

[The official jade language reference](http://jade-lang.com/reference/)

### Html Templates

You may use go's [html/template](http://golang.org/pkg/html/template/) package for html templates and layouts.
Admittedly, this may be the component that is the hardest to grok at first. I strongly recommend the related
resources at the bottom of this section.

All templates should have the .tmpl extension. Scribble uses the following conventions, but you are not 
strictly required to adhere to them:

1. Html layouts (which exist in `layoutsDir`) should be a wrapper around a named template, typically called
`"content"`. Templates which use the layout define a content named template and then render the layout with
the content inside. 
2. Layouts and includes are identified by their base filename, including the extension. (This is the go
default when using `ParseFiles` or `ParseGlob`).
3. Post layouts define a content template and may use a layout as a wrapper.

Here are some examples of converting layouts from the [seed project](https://github.com/albrow/scribble-seed)
to go's native templates:

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
<!-- Render the content named template inside of the base.tmpl layout -->
{{ template "base.tmpl" }}
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
<!-- Render the content named template inside of the base.tmpl layout -->
{{ template "base.tmpl" . }}
```

#### Related Resources:

- [Learn more about go's text templates](http://golang.org/pkg/text/template/), which share a lot of functionality with html templates.
- [Learn more about go's html templates](http://golang.org/pkg/html/template/).
- [Learn more about layout/template inheritance in go](https://elithrar.github.io/article/approximating-html-template-inheritance/). (The ideas there will work well with scribble).


License
-------

Scribble is licensed under the MIT License. See the LICENSE file for more information.
