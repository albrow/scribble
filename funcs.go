package main

import (
	"html/template"
)

var funcMap template.FuncMap = map[string]interface{}{
	// Posts returns up to limit posts. If limit is 0, it returns
	// all posts. If limit is greater than len(posts), it returns
	// all posts.
	"Posts": func(limit ...int) []*Post {
		if len(limit) == 0 || limit[0] == 0 || limit[0] > len(posts) {
			return posts
		} else {
			return posts[:limit[0]]
		}
	},
}
