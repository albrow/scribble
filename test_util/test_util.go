package test_util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// CheckStringsMatch adds an error to t iff the elements in got do
// not match exactly the elements in expected (irrespective of order).
func CheckStringsMatch(t *testing.T, expected []string, got []string) {
	gotMap, expectedMap := map[string]struct{}{}, map[string]struct{}{}
	for _, s := range got {
		gotMap[s] = struct{}{}
	}
	for _, s := range expected {
		expectedMap[s] = struct{}{}
	}
	if !reflect.DeepEqual(gotMap, expectedMap) {
		t.Errorf("Paths were not correct.\nExpected: %v\nGot: %v\n", expected, got)
	}
}

// CheckFilesMatch adds an error to t iff the contents of the file at gotPath
// do not match exactly the contents of the file at expectedPath.
func CheckFilesMatch(t *testing.T, expectedPath string, gotPath string) {
	expected, err := ioutil.ReadFile(expectedPath)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := ioutil.ReadFile(gotPath); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("File at %s did not exist.", gotPath)
		} else {
			t.Fatal(err)
		}
	} else if len(got) == 0 {
		t.Errorf("File at %s was empty.", gotPath)
	} else if !reflect.DeepEqual(expected, got) {
		t.Errorf("Contents of file at %s were incorrect.\nExpected: %s\nGot: %s\n", gotPath, string(expected), string(got))
	}
}

// CheckDirsMatch recursively iterates through expectedDir and checks that the directory
// structure and the contents of each file match gotDir exactly. If anything does not
// match, it adds an error to t.
func CheckDirsMatch(t *testing.T, expectedDir string, gotDir string) {
	if err := filepath.Walk(expectedDir, func(expectedPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !(expectedPath == expectedDir) {
			// We expect the directory structure to be the same, so every subdirectory
			// following expectedDir should also be present in gotDir. i.e. if expectedDir is
			// /tmp/source, gotDir is /tmp/public, and expectedPath is
			// /tmp/source/one/two/three.txt, we would expect the corresponding gotPath to
			// be /tmp/public/one/two/three.txt.
			gotPath := strings.Replace(expectedPath, expectedDir, gotDir, 1)
			CheckFilesMatch(t, expectedPath, gotPath)
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
