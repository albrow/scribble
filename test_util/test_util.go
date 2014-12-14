package test_util

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// CheckStringsMatch adds an error to t iff the elements in got do
// not match exactly the elements in expected (irrespective of order).
func CheckStringsMatch(t *testing.T, expected []string, got []string) {
	gotMap, expectedMap := map[string]struct{}{}, map[string]struct{}{}
	for _, path := range got {
		gotMap[path] = struct{}{}
	}
	for _, path := range expected {
		expectedMap[path] = struct{}{}
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
