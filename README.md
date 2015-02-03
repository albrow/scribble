Scribble
========

Version: X.X.X

A tiny static blog generator written in go.

Although I am still planning to add some awesome features before v1.0, the core functionality
for scribble is finished. This package can be considered safe and stable enough for general use.


How it Works
------------

Scribble is a command line tool that compiles different types of source files into a
static blog made up of html, css, and js. It uses:

- Markdown for writting posts
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
`git clone https://github.com/albrow/scribble-seed <optional-blog-name-here>`.

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

`config.toml` stores some basic configuration in [toml](https://github.com/toml-lang/toml) format.
This consists of metadata such as your Blog's title, the author's name, and a description, but also
tells scribble where to look for certain files. Every scribble project must have a config.toml file
in the project root directory.

`public` is the folder where scribble will put your finished website after compiling. It's also the
folder that scribble will serve from when using the `scribble serve` command. This is set via the
`destDir` key in `config.toml`.

`source` is where all the source code you write will live. This includes things like stylesheets,
posts, html templates, and javascript. This is set via the `sourceDir` key in `config.toml`.

`source/_includes` is an optional folder where you can put partial html templates, i.e. templates
which don't constitute a full page on their own, but are meant to be *included* in other templates.
In the seed project, we use there are two files in `source/_includes`, one for filling in the `<head>`
tag with metadata and stylesheets, and one for including any javascript files at the bottom of the
`<body>` tag. This is set via the `includesDir` key in `config.toml`.

`source/_layouts` is where you put html layouts. Layouts are reusable html wrappers that define how
each pages will look. In the seed project, there is just one layout, called `base.tmpl`. It consists
of html boilerplate like the `<html>`, `<head>`, and `<body>` tags. It includes the two files in our
`_includes` directory by using the template/html syntax `{{ template "head.tmpl" . }}`. The layouts
directory is defined via the `layoutsDir` key in `config.toml`.

`source/_post_layouts` is where you put post layouts. Like html layouts, post layouts are reusable
html wrappers that define each post will look. In the seed project, there is just one post layout,
called `post.tmpl`. The default template simply consists of the title of the post in a header and
the content of the post wrapped in a div below it. The post layouts directory is defined via the
`postLayoutsDir` key in `config.toml`.

`_posts` is where your posts will reside. Posts are written in markdown and must include toml
frontmatter which defines the post layout, and optionally the title, author's name, date, and
description. The posts directory is defined via the `postsDir` key in `config.toml`.

`index.tmpl` is the index page and will compile to index.html. It consists of an unordered list
of links to the 5 most recent posts.

`js` is a folder where you can put javascript. Any files here will be copied to `destDir` directly.

`styles` is where you can put sass stylesheets. They will be picked up by the sass compiler, and
any sass files that don't start with an underscore will be compiled to css files in `destDir`.


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

`index.html` came from `source/index.tmpl`

`js/main.js` came from `source/js/main.js`

`one/index.html` came from `source/_posts/one.md`, and the other posts - `two` and `three` came from their
respective source files.

`styles/main.css` came from `source/styles/main.scss`

Note that the `_includes`, `_layouts`, and `_post_layouts` folders were not compiled because they started
with underscores. Same with the imported sass files: `_fonts.scss` and `_colors.scss`.

### Posts

TODO: Fill this out

### Sass

TODO: Fill this out

### Html Templates

TODO: Fill this out
