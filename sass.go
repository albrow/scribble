package main

import (
	"fmt"
	"os"
	"os/exec"
)

func compileSass(watch bool) {
	// TODO: implement watch
	// TODO: switch to sassc
	fmt.Println("    compiling sass")

	// set up and execute the command, piping output to stdout
	cmd := exec.Command("sass", "--update", fmt.Sprintf("%s:%s", sassSourceDir, sassDestDir))
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
