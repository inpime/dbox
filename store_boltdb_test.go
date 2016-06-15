package dbox

import (
	"github.com/boltdb/bolt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func mustOpenBoltDB(file string) *bolt.DB {
	if err := ensureDir(filepath.Dir(file)); err != nil {
		panic(err)
	}

	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		panic(err)
	}

	return db
}

func TestBoltDBStore_simpleStrategy(t *testing.T) {
	storePath := "./boltdb.db"
	store := NewBoltDBStore(mustOpenBoltDB("./boltdb.db"), "bucketname")

	createSimpleStrategy(t, store, store, store)

	deleteSimpleStrategy(t, store, store, store)

	os.RemoveAll(storePath)
}

func BenchmarkTestBoltDBStore_simpleStrategy(b *testing.B) {
	storePath := "./boltdb.db"
	store := NewBoltDBStore(mustOpenBoltDB("./boltdb.db"), "bucketname")

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
	store := NewBoltDBStore(mustOpenBoltDB("./boltdb.db"), "bucketname")

	file := NewFile(store)
	mapSet(file.Meta(), "a", "b")

	file.RawData().Write([]byte("text text"))
	mapSet(file.MapData(), "map1", "v1")
	file.Sync()

	fileId := file.ID()

	for i := 0; i < b.N; i++ {
		_file, _ := NewFileID(fileId, store)
		_file.ID()
		_file.MapData()
		_file.RawData()
	}

	file.Delete()

	os.RemoveAll(storePath)
}
