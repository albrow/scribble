package main

import (
	"fmt"
	"github.com/albrow/scribble/config"
	"github.com/codegangsta/negroni"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// serve serves all the static content in config.DestDir via a lightweight
// negroni server on the given port.
func serve(port int) {
	fmt.Printf("--> serving on port %d\n", port)
	// use negroni to serve destDir on port
	destFileSystem := http.Dir(config.DestDir)
	n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(destFileSystem), negroni.HandlerFunc(NotFound))
	n.Run(fmt.Sprintf(":%d", port))
}

func NotFound(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Header().Add("Content-Type", "text/html")
	urlPath := strings.Replace(r.URL.String(), "/", string(os.PathSeparator), -1)
	lookedPath := filepath.Join(config.DestDir, urlPath)
	content := fmt.Sprintf("<h3>404 Not Found</h3><p>Scribble could not find <em>%s</em>. Looked in <em>%s</em>.</p>", r.URL, lookedPath)
	fmt.Fprint(rw, wrapHtml("Not Found", content))
}

// wrapHtml returns a string of boilerplate-wrapped html with the given title and content.
func wrapHtml(title string, content string) string {
	return fmt.Sprintf(`<!doctype html><html><head><title>%s</title></head><body>%s</body></html>`, title, content)
}
