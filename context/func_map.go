package context

import (
	"html/template"
)

// FuncMap represents a set of functions, identified by some key, which
// will be availalbe to templates when rendering. Similarly to the context,
// the FuncMap will be passed through any time an ace template is rendered.
// See http://golang.org/pkg/text/template/#FuncMap. There is one func provided
// by default, called Posts, which is defined in generators/posts_compiler.
var FuncMap template.FuncMap = map[string]interface{}{}
