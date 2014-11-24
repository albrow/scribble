package lib

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

// a list of config vars
var (
	SourceDir, DestDir, PostsDir, LayoutsDir string
)

// parseConfig reads and parses config.toml, setting the values
// of the above config variables.
func parseConfig() {
	fmt.Println("--> parsing config.toml")
	if _, err := toml.DecodeFile("config.toml", context); err != nil {
		panic(err)
	}
	vars := map[string]*string{
		"sourceDir":  &SourceDir,
		"destDir":    &DestDir,
		"postsDir":   &PostsDir,
		"layoutsDir": &LayoutsDir,
	}
	setGlobalConfig(vars, context)
}

// setGlobalConfig sets the values of vars based on the contents of data
func setGlobalConfig(vars map[string]*string, data map[string]interface{}) {
	for name, holder := range vars {
		if value, found := data[name]; found {
			(*holder) = fmt.Sprint(value)
		}
	}
}
