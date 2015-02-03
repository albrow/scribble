package main

import (
	"fmt"
	"github.com/albrow/scribble/util"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
)

var (
	app = kingpin.New("scribble", "A tiny static blog generator written in go.")

	versionCmd = app.Command("version", "Display version information and then quit.")

	serveCmd   = app.Command("serve", "Compile and serve the site.")
	servePort  = serveCmd.Flag("port", "The port on which to serve the site.").Short('p').Default("4000").Int()
	serveTrace = serveCmd.Flag("trace", "Whether or not to print a full stack trace when there is an error.").Short('t').Default("false").Bool()

	compileCmd   = app.Command("compile", "Compile the site.")
	compileWatch = compileCmd.Flag("watch", "Whether or not to watch for changes and automatically recompile.").Short('w').Default("").Bool()
	compileTrace = compileCmd.Flag("trace", "Whether or not to print a full stack trace when there is an error.").Short('t').Default("false").Bool()
)

const (
	version = "0.3.0"
)

func main() {
	// catch panics and print them out as errors
	defer util.Recovery(*compileTrace || *serveTrace)

	// Parse the command line arguments and flags and delegate
	// to the appropriate functions.
	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Usage(os.Stdout)
		os.Exit(0)
	}
	switch cmd {
	case versionCmd.FullCommand():
		fmt.Println("scribble version:", version)
	case compileCmd.FullCommand():
		compile(*compileWatch)
		if *compileWatch {
			// If the watch flag was provided, don't exit.
			// User will need to quit manually, e.g. with ctrl+c
			done := make(chan bool)
			<-done
		}
	case serveCmd.FullCommand():
		compile(true)
		serve(*servePort)
	default:
		app.Usage(os.Stdout)
		os.Exit(0)
	}
}
