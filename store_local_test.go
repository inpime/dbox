package dbox

import (
	"io/ioutil"
	"os"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestLocalStore_simpleStrategy(t *testing.T) {
	store := NewLocalStore("/Users/gbv/work/src/dbox/_testdata/")

	createSimpleStrategy(t, store)

	deleteSimpleStrategy(t, store)
}

func TestLocalStore_writereadFile(t *testing.T) {
	ioutil.WriteFile("testfile", []byte("text"), 0644)
	ioutil.WriteFile("testfile", []byte("text text"), 0644)
	b, _ := ioutil.ReadFile("testfile")
	if string(b) != "text text" {
		t.Error("not expected value")
	}
	os.Remove("testfile")
}
