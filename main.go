package main

import (
	"github.com/albrow/scribble/lib"
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

func main() {
	defer lib.Recovery()
	kingpin.Version(version)
	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Usage(os.Stdout)
		os.Exit(0)
	}
	switch cmd {
	case compileCmd.FullCommand():
		lib.Compile(*compileWatch)
	case serveCmd.FullCommand():
		lib.Compile(true)
		lib.Serve(*servePort)
	default:
		app.Usage(os.Stdout)
		os.Exit(0)
	}
}
