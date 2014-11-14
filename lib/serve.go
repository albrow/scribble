package lib

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
)

func Serve(port int) {
	fmt.Printf("--> serving on port %d\n", port)
	// use negroni to serve destDir on port
	destFileSystem := http.Dir(destDir)
	n := negroni.New(negroni.NewStatic(destFileSystem), negroni.NewRecovery())
	n.Run(fmt.Sprintf(":%d", port))
}
