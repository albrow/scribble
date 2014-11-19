Scribble
========

A tiny static blog generator written in go.

This isn't quite ready yet. I'm still planning to add some pretty crucial features,
clean up the code, and make it a little simpler to use. I will update this README
with more details after I make those changes.

In the meantime, if you're looking for a static blog generator written in go, check out
[Hugo](https://github.com/spf13/hugo).


How it Works
------------

Scribble is a command line tool that compiles different types of source files into a
static blog made up of html, css, and js. It uses:

- Markdown for writting posts
- Toml for frontmatter and configuration
- Sass for styling
- Ace for html templates.

Scribble is highly optimized for speed. It uses sassc (a C port of the sass compiler). It
also features watching/livereload and only recompiles the necessary files. The result of this
speed is a really pleasant design/development experience. For exapmle, if you make a change to
a sass file, scribble will only recompile the affected sass file(s) and LiveReload will only
reload the changed css file(s) without refreshing the entire page. That means you can get instant
feedback. (Especially useful when styling dropdowns and other transient elements which would
dissappear if the entire page reloaded!)

Why?
----

I didn't build Scribble because I thought the world needed another static blog generator
(there are some fantastic ones out there already!). I built it because I wanted practice
working with generated files with go, and because I wanted to use my favorite markup
languages and preprocessors (sass, ace, toml, markdown). I built it for myself. Because
of this, Scribble is pretty opinionated and might not make sense for everyone.
