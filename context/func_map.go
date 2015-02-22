// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package context

import (
	"html/template"
)

// FuncMap represents a set of functions, identified by some key, which
// will be availalbe to templates when rendering. Similarly to the context,
// the FuncMap will be passed through any time an ace template is rendered.
// See http://golang.org/pkg/text/template/#FuncMap. There is one func provided
// by default, called Posts, which is defined in compilers/posts_compiler.
var FuncMap template.FuncMap = map[string]interface{}{}
