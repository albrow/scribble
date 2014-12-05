package lib

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestCreateFile(t *testing.T) {
	defer func() {
		// recover from any potential panics and make the
		// test fail if there were any
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()
	defer func() {
		// cleanup by removing all the files we created
		os.RemoveAll("/tmp/test")
	}()
	rand.Seed(time.Now().Unix())
	path := fmt.Sprintf("/tmp/test/testFile%d.txt", rand.Int())
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
