package main

import (
	"github.com/albrow/scribble/util"
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
	version = "X.X.X (develop)"
)

func main() {
	// catch panics and print them out as errors
	defer util.Recovery()
	// print out the version when prompted
	kingpin.Version(version)

	// Parse the command line arguments and flags and delegate
	// to the appropriate functions.
	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Usage(os.Stdout)
		os.Exit(0)
	}
	switch cmd {
	case compileCmd.FullCommand():
		compile(*compileWatch)
	case serveCmd.FullCommand():
		compile(true)
		serve(*servePort)
	default:
		app.Usage(os.Stdout)
		os.Exit(0)
	}
}
