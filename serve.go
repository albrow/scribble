package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
)

func serve(port int) {
	fmt.Println("--> serving")
	fmt.Printf("    port: %v\n", port)

	// use negroni to serve destDir on port
	destFileSystem := http.Dir(destDir)
	n := negroni.New(negroni.NewStatic(destFileSystem), negroni.NewRecovery())
	n.Run(fmt.Sprintf(":%d", port))
}
