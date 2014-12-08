package generators

import (
	"github.com/albrow/scribble/util"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func createEmptyFiles(paths []string) error {
	for _, path := range paths {
		if f, err := util.CreateFileWithPath(path); err != nil {
			return err
		} else {
			f.Close()
		}
	}
	return nil
}

func checkPathsMatch(t *testing.T, compiler Compiler, srcDir string, expectedPaths []string) {
	gotPaths, err := FindPaths(srcDir, compiler)
	if err != nil {
		t.Fatal(err)
	}
	gotPathMap, expectedPathMap := map[string]struct{}{}, map[string]struct{}{}
	for _, path := range gotPaths {
		gotPathMap[path] = struct{}{}
	}
	for _, path := range expectedPaths {
		expectedPathMap[path] = struct{}{}
	}
	if !reflect.DeepEqual(gotPathMap, expectedPathMap) {
		t.Errorf("Paths were not correct.\nExpected: %v\nGot: %v\n", expectedPaths, gotPaths)
	}
}

// expects a map of destination path to source path
func copyFiles(paths map[string]string) error {
	for destPath, srcPath := range paths {
		dest, err := util.CreateFileWithPath(destPath)
		if err != nil {
			return err
		}
		src, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dest, src); err != nil {
			return err
		}
	}
	return nil
}

func checkOutputMatchesFile(t *testing.T, expectedPath string, gotPath string) {
	expected, err := ioutil.ReadFile(expectedPath)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := ioutil.ReadFile(gotPath); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Expected compiled result at %s but file did not exist.", gotPath)
		} else {
			t.Fatal(err)
		}
	} else if len(got) == 0 {
		t.Errorf("Compiled result at %s was empty.", gotPath)
	} else if !reflect.DeepEqual(expected, got) {
		t.Errorf("Compiled result at %s was incorrect.\nExpected: %s\nGot: %s\n", gotPath, string(expected), string(got))
	}
}
