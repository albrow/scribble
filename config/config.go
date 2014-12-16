package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/scribble/context"
)

// a list of config vars
var (
	SourceDir, DestDir, PostsDir, LayoutsDir, ViewsDir string
)

// Parse reads and parses config.toml, setting the values
// of the config variables here and in the context.
func Parse() {
	fmt.Println("--> parsing config.toml")
	if _, err := toml.DecodeFile("config.toml", context.GetContext()); err != nil {
		msg := fmt.Sprintf("Problem reading config.toml file:\n%s", err)
		panic(msg)
	}
	vars := map[string]*string{
		"sourceDir":  &SourceDir,
		"destDir":    &DestDir,
		"postsDir":   &PostsDir,
		"layoutsDir": &LayoutsDir,
		"viewsDir":   &ViewsDir,
	}
	setConfig(vars, context.GetContext())
}

// setConfig sets the values of vars based on the contents of data
func setConfig(vars map[string]*string, data map[string]interface{}) {
	for name, holder := range vars {
		if value, found := data[name]; found {
			(*holder) = fmt.Sprint(value)
		}
	}
}
