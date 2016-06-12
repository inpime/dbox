package dbox

import (
	"os"
	"testing"
)

func TestLocalStore_simpleStrategy(t *testing.T) {
	store := NewLocalStore("./_testdata/")

	createSimpleStrategy(t, store)

	deleteSimpleStrategy(t, store)

	os.RemoveAll("./_testdata/")
}
