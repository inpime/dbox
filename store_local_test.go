package dbox

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStore_simpleStrategy(t *testing.T) {
	storePath := "./_testdata/"
	store := NewLocalStore(storePath)

	fileId := createSimpleStrategy(t, store, store, store)

	file := NewFile(store)
	file.SetMapDataStore(store)
	file.SetRawDataStore(store)

	err := store.Get(fileId, file)
	assert.NoError(t, err, "get file by id for checker")

	// check exist files
	assert.True(t, exists(storePath+file.ID()), "check exist file")
	assert.True(t, exists(storePath+file.Name()), "check exist file")
	file.MapData() // force init map object
	assert.True(t, exists(storePath+file.mdata.ID()), "check exist file")
	file.RawData() // force init raw object
	assert.True(t, exists(storePath+file.rdata.ID()), "check exist file")

	// check count files
	count := 0
	err = filepath.Walk(storePath, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			count++
		}
		return nil
	})
	assert.NoError(t, err, "wall files")
	assert.Equal(t, count, 4, "count files")

	deleteSimpleStrategy(t, store, store, store)

	count = 0
	err = filepath.Walk(storePath, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			count++
		}
		return nil
	})
	assert.NoError(t, err, "wall files")
	assert.Equal(t, count, 0, "count files")

	os.RemoveAll(storePath)
}

func BenchmarkLocalFileSystem_simpleStrategy(b *testing.B) {
	storePath := "./_testdata/"
	store := NewLocalStore(storePath)

	for i := 0; i < b.N; i++ {
		file := NewFile(store)
		mapSet(file.Meta(), "a", "b")

		file.RawData().Write([]byte("text text"))
		mapSet(file.MapData(), "map1", "v1")
		file.Sync()

		fileId := file.ID()

		// Load

		file = NewFile(store)

		store.Get(fileId, file)

		file.Delete()
	}
}
