package main

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/albrow/ace"
	"github.com/wsxiaoys/terminal/color"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func compilePages() {
	fmt.Println("--> compiling pages")
	// walk through the source dir
	if err := filepath.Walk(sourceDir, func(innerPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		base := info.Name()
		if base[0] == '.' || base[0] == '_' {
			// ignore two kinds of files
			// 1. those that start with a '.' are hidden system files
			// 2. those that start with a '_' are specifically ignored by scribble
			if info.IsDir() {
				// skip any files in directories that start with '_'
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(base)
			switch ext {
			case ".ace":
				compilePageFromPath(innerPath)
			default:
				// copy the file directly to the destDir
				destPath := strings.Replace(innerPath, sourceDir, destDir, 1)
				srcFile, err := os.Open(innerPath)
				if err != nil {
					panic(err)
				}
				if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
					panic(err)
				}
				destFile, err := os.Create(destPath)
				if err != nil {
					panic(err)
				}
				if _, err := io.Copy(destFile, srcFile); err != nil {
					panic(err)
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func compilePageFromPath(path string) {
	srcFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(srcFile)
	frontMatter, content, err := split(reader)
	pageContext := context
	if frontMatter != "" {
		if _, err := toml.Decode(frontMatter, pageContext); err != nil {
			chimeError(err)
		}
	}
	layout := "base"
	if otherLayout, found := pageContext["layout"]; found {
		layout = otherLayout.(string)
	}
	tpl, err := ace.Load("_layouts/"+layout, filepath.Base(path), &ace.Options{
		DynamicReload: true,
		BaseDir:       sourceDir,
		FuncMap:       funcMap,
		Asset: func(name string) ([]byte, error) {
			return []byte(content), nil
		},
	})
	if err != nil {
		chimeError(err)
	}
	destPath := strings.Replace(path, sourceDir, destDir, 1)
	destPath = strings.Replace(destPath, ".ace", ".html", 1)
	if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
		panic(err)
	}
	color.Printf("@g    CREATE: %s -> %s\n", path, destPath)
	destFile, err := os.Create(destPath)
	if err != nil {
		// if the file already exists, that's fine
		// if there was some other error, panic
		if !os.IsExist(err) {
			panic(err)
		}
	}
	if err := tpl.Execute(destFile, pageContext); err != nil {
		chimeError(err)
	}
}
