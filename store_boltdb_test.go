package dbox

import (
	"os"
	"testing"
)

func TestBoltDBStore_simpleStrategy(t *testing.T) {
	storePath := "./boltdb.db"
	store := NewBoltDBStore(storePath)

	createSimpleStrategy(t, store, store, store)

	deleteSimpleStrategy(t, store, store, store)

	os.RemoveAll(storePath)
}

func BenchmarkTestBoltDBStore_simpleStrategy(b *testing.B) {
	storePath := "./boltdb.db"
	store := NewBoltDBStore(storePath)

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

	os.RemoveAll(storePath)
}

func BenchmarkTestBoltDBStore_onlyReadFile(b *testing.B) {
	storePath := "./boltdb.db"
	store := NewBoltDBStore(storePath)

	file := NewFile(store)
	mapSet(file.Meta(), "a", "b")

	file.RawData().Write([]byte("text text"))
	mapSet(file.MapData(), "map1", "v1")
	file.Sync()

	fileId := file.ID()

	for i := 0; i < b.N; i++ {

		file = NewFile(store)

		store.Get(fileId, file)
	}

	file.Delete()

	os.RemoveAll(storePath)
}
