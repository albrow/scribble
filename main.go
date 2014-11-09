package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
)

var (
	app = kingpin.New("scribble", "A tiny static blog generator written in go.")

	serveCmd  = app.Command("serve", "Compile and serve the site.")
	servePort = serveCmd.Flag("port", "The port on which to serve the site.").Short('p').Default("4000").Int()

	compileCmd   = app.Command("compile", "Compile the site.")
	compileWatch = compileCmd.Flag("watch", "Whether or not to watch for changes and automatically recompile.").Short('w').Default("").Bool()
)

const (
	version = "0.0.1"
)

type Context map[string]interface{}

var context = Context{}

func main() {
	kingpin.Version(version)
	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Usage(os.Stdout)
		os.Exit(0)
	}
	switch cmd {
	case compileCmd.FullCommand():
		compile(*compileWatch)
	case serveCmd.FullCommand():
		serve(*servePort)
	default:
		app.Usage(os.Stdout)
		os.Exit(0)
	}
}

func serve(port int) {
	compile(true)
	fmt.Println("--> serving")
	fmt.Printf("    port: %v\n", port)
}
