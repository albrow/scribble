package main

import (
	"fmt"
	"os"
	"os/exec"
)

func compileSass(watch bool) {
	fmt.Println("    compiling sass")
	// choose the appropriate flag depending on the value of watch
	var flag = ""
	if watch {
		flag = "--watch"
	} else {
		flag = "--update"
	}

	// set up and execute the command, piping output to stdout
	cmd := exec.Command("sass", flag, fmt.Sprintf("%s:%s", sassSourceDir, sassDestDir))
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
