package util

import (
	"fmt"
	"github.com/albrow/scribble/log"
)

// ChimeError outputs the bell character and then the error message,
// colored red and formatted.
func ChimeError(err interface{}) {
	fmt.Print("\a")
	log.Error.Printf("ERROR: %s\n", err)
}

// ChimeErrorf outputs the bell character and then the error message,
// colored red and formatted according to format and args. It works
// just like fmt.Printf.
func ChimeErrorf(format string, args ...interface{}) {
	fmt.Print("\a")
	log.Error.Printf("ERROR: %s\n", fmt.Sprintf(format, args...))
}
