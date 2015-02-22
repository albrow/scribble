// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package main

import (
	"github.com/albrow/scribble/compilers"
	"github.com/albrow/scribble/config"
	"github.com/albrow/scribble/log"
	"github.com/albrow/scribble/util"
	"os"
)

// compile compiles all the contents of config.SourceDir and puts the compiled
// result in config.DestDir.
func compile(watch bool) {
	config.Parse()
	log.Default.Println("Compiling...")
	if err := createDestDir(); err != nil {
		panic(err)
	}
	if err := compilers.CompileAll(); err != nil {
		util.ChimeError(err)
	}
	if watch {
		watchAll()
	}
}

func createDestDir() error {
	if err := os.MkdirAll(config.DestDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			// If the directory already existed, that's fine,
			// otherwise return the error
			return err
		}
	}
	return nil
}
