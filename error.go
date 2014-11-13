package main

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
)

// chimeError outputs the bell character and then the error message,
// colored red and formatted.
func chimeError(err interface{}) {
	fmt.Print("\a")
	color.Printf("@r    ERROR: %s\n", err)
}

func chimeErrorf(format string, args ...interface{}) {
	fmt.Print("\a")
	color.Printf("@r    ERROR: %s\n", fmt.Sprintf(format, args...))
}
