package util

import (
	"github.com/albrow/scribble/test_util"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateFileWithPath(t *testing.T) {
	root := "/tmp/test_create_file"
	defer func() {
		// cleanup by removing all the files we created
		os.RemoveAll(root)
	}()
	path := root + "/testFile.txt"
	// Open a temp file and write to it
	f, err := CreateFileWithPath(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	data := "Hello, this is a test!\n"
	if _, err := f.WriteString(data); err != nil {
		t.Error(err)
	}
	if err := f.Sync(); err != nil {
		t.Error(err)
	}
	// Open the same temp file and read from it
	gotData, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	if string(gotData) != data {
		t.Errorf("Read data was not correct. Expected %s but got %s.\n", data, string(gotData))
	}
}

func TestCopyFile(t *testing.T) {
	t.Skip("TODO: Finish TestCopyFile")
	root := "/tmp/test_copy_files"
	defer func() {
		// cleanup by removing all the files we created
		os.RemoveAll(root)
	}()
	// Create and write to a test file
	srcPath := root + "/original/testFile.txt"
	f, err := CreateFileWithPath(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString("Hello, this is a test!\n")
	// Copy the test file
	destPath := root + "/copy/testFile.txt"
	if err := CopyFile(srcPath, destPath); err != nil {
		t.Fatal(err)
	}
	test_util.CheckFilesMatch(t, srcPath, destPath)
}

func TestRecursiveCopy(t *testing.T) {
	t.Skip("TODO: Finish TestRecursiveCopy")
	destDir := "/tmp/test_recursive_copy"
	defer func() {
		// cleanup by removing all the files we created
		os.RemoveAll(destDir)
	}()
	srcDir := os.Getenv("GOPATH") + "src/github.com/albrow/scribble/test_files/sass/source"
	if err := RecursiveCopy(srcDir, destDir); err != nil {
		t.Fatal(err)
	}
	test_util.CheckDirsMatch(t, srcDir, destDir)
}
