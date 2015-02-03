package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/scribble/context"
	"github.com/albrow/scribble/log"
)

// a list of config vars
var (
	SourceDir, DestDir, PostsDir, LayoutsDir, PostLayoutsDir, IncludesDir string
)

// Parse reads and parses config.toml, setting the values
// of the config variables here and in the context. It panics
// if there was a problem reading the file.
func Parse() {
	log.Default.Println("Parsing config.toml...")
	if _, err := toml.DecodeFile("config.toml", context.GetContext()); err != nil {
		msg := fmt.Sprintf("Problem reading config.toml file:\n%s", err)
		panic(msg)
	}
	vars := map[string]*string{
		"sourceDir":      &SourceDir,
		"destDir":        &DestDir,
		"layoutsDir":     &LayoutsDir,
		"postsDir":       &PostsDir,
		"postLayoutsDir": &PostLayoutsDir,
		"includesDir":    &IncludesDir,
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
