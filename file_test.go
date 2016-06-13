package dbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createSimpleStrategy(t *testing.T, store, mapStore, rawStore Store) string {
	file := NewFile(store)
	file.SetName("namefile")
	file.SetMapDataStore(mapStore)
	file.SetRawDataStore(rawStore)

	var err error
	err = file.Sync()

	mapSet(file.Meta(), "a", "b")
	// err = file.Sync()

	assert.NoError(t, err, "sync file")

	mapSet(file.Meta(), "d", "b")

	mapSet(file.MapData(), "map1", "v1")

	// err = file.Sync()
	assert.NoError(t, err, "sync file")

	file.RawData().Write([]byte("text text"))
	// err = file.Sync()
	assert.NoError(t, err, "sync file")

	mapSet(file.MapData(), "map2", "v2")

	err = file.Sync()
	assert.NoError(t, err, "sync file")

	fileId := file.ID()

	// Load by id

	file, err = NewFileID(fileId, store)
	file.SetMapDataStore(mapStore)
	file.SetRawDataStore(rawStore)

	assert.NoError(t, err, "get by id %q", fileId)

	assert.Equal(t, mapString(file.Meta(), "a"), "b", "not expected value")
	assert.Equal(t, mapString(file.Meta(), "d"), "b", "not expected value")
	assert.Equal(t, mapString(file.MapData(), "map1"), "v1", "not expected value")
	assert.Equal(t, mapString(file.MapData(), "map2"), "v2", "not expected value")
	assert.Equal(t, file.RawData().Bytes(), []byte("text text"), "not expected value")

	// Load by name
	file, err = NewFileName("namefile", store)
	file.SetMapDataStore(mapStore)
	file.SetRawDataStore(rawStore)
	assert.NoError(t, err, "get by name")

	assert.Equal(t, mapString(file.Meta(), "a"), "b", "not expected value")
	assert.Equal(t, mapString(file.Meta(), "d"), "b", "not expected value")
	assert.Equal(t, mapString(file.MapData(), "map1"), "v1", "not expected value")
	assert.Equal(t, mapString(file.MapData(), "map2"), "v2", "not expected value")
	assert.Equal(t, file.RawData().Bytes(), []byte("text text"), "not expected value")

	return fileId
}

func deleteSimpleStrategy(t *testing.T, store, mapStore, rawStore Store) string {
	file, err := NewFileName("namefile", store)
	file.SetMapDataStore(mapStore)
	file.SetRawDataStore(rawStore)
	assert.NoError(t, err, "get by name")

	fileId := file.ID()

	err = file.Delete()
	assert.NoError(t, err, "remove file")

	// Check remove

	file, err = NewFileID(fileId, store)
	file.SetMapDataStore(mapStore)
	file.SetRawDataStore(rawStore)
	assert.Equal(t, err, ErrNotFound, "get removed file by id")

	file, err = NewFileName("namefile", store)
	file.SetMapDataStore(mapStore)
	file.SetRawDataStore(rawStore)
	assert.Equal(t, err, ErrNotFound, "get removed file by nmae")

	return fileId
}

func TestFile_simpleStrategy(t *testing.T) {
	store := NewMemoryStore()

	createSimpleStrategy(t, store, store, store)
	assert.Equal(t, len((*MemoryStore)(store).list), 4, "not valid storage")

	deleteSimpleStrategy(t, store, store, store)
	assert.Equal(t, len((*MemoryStore)(store).list), 0, "not valid storage")
}

func BenchmarkFile_simpleStrategy(b *testing.B) {
	store := NewMemoryStore()

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
