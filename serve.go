package main

import (
	"fmt"
	"github.com/albrow/scribble/config"
	"github.com/codegangsta/negroni"
	"net/http"
)

// serve serves all the static content in config.DestDir via a lightweight
// negroni server on the given port.
func serve(port int) {
	fmt.Printf("--> serving on port %d\n", port)
	// use negroni to serve destDir on port
	destFileSystem := http.Dir(config.DestDir)
	n := negroni.New(negroni.NewStatic(destFileSystem), negroni.NewRecovery())
	n.Run(fmt.Sprintf(":%d", port))
}
