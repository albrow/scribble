package lib

import (
	"html/template"
)

var FuncMap template.FuncMap = map[string]interface{}{
	"Posts": Posts,
}
