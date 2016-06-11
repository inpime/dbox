package dbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createSimpleStrategy(t *testing.T, store Store) {
	file := NewFile(store)
    file.SetName("namefile")

	file.Meta().Set("a", "b")
	err := file.Sync()
    assert.NoError(t, err, "get by id")

	file.Meta().Set("d", "b")

	file.RawData().Write([]byte("text text"))
	err = file.Sync()
    assert.NoError(t, err, "get by id")

	file.MapData().Set("map1", "v1")
	file.MapData().Set("map2", "v2")

	err = file.Sync()
    assert.NoError(t, err, "get by id")

	fileId := file.ID()

	// Load by id

	file, err = NewFileID(fileId, store)
	assert.NoError(t, err, "get by id")

	assert.Equal(t, file.Meta().String("a"), "b", "not expected value")
	assert.Equal(t, file.Meta().String("d"), "b", "not expected value")
	assert.Equal(t, file.MapData().String("map1"), "v1", "not expected value")
	assert.Equal(t, file.MapData().String("map2"), "v2", "not expected value")
	assert.Equal(t, file.RawData().Bytes(), []byte("text text"), "not expected value")

    // Load by name
    file, err = NewFileName("namefile", store)
	assert.NoError(t, err, "get by id")

	assert.Equal(t, file.Meta().String("a"), "b", "not expected value")
	assert.Equal(t, file.Meta().String("d"), "b", "not expected value")
	assert.Equal(t, file.MapData().String("map1"), "v1", "not expected value")
	assert.Equal(t, file.MapData().String("map2"), "v2", "not expected value")
	assert.Equal(t, file.RawData().Bytes(), []byte("text text"), "not expected value")
}

func deleteSimpleStrategy(t *testing.T, store Store) {
	file, err := NewFileName("namefile", store)
	assert.NoError(t, err, "get by id")

	fileId := file.ID()

	err = file.Delete() 
    assert.NoError(t, err, "remove file")

    // Check remove

    file, err = NewFileID(fileId, store)
	assert.Equal(t, err, ErrNotFound, "get removed file by id")

    file, err = NewFileName("namefile", store)
	assert.Equal(t, err, ErrNotFound, "get removed file by nmae")
}

func TestFile_simpleStrategy(t *testing.T) {
	store := NewMemoryStore()

	createSimpleStrategy(t, store)
	assert.Equal(t, len((*MemoryStore)(store).list), 4, "not valid storage")

	deleteSimpleStrategy(t, store)
    assert.Equal(t, len((*MemoryStore)(store).list), 0, "not valid storage")
}

func BenchmarkFile_simpleStrategy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		store := NewMemoryStore()
		file := NewFile(store)
		file.Meta().Set("a", "b")

        file.RawData().Write([]byte("text text"))
        file.MapData().Set("map1", "v1")
		file.Sync()
		
		fileId := file.ID()

		// Load

		file = NewFile(store)

		store.Get(fileId, file)

        file.Delete() 
	}
}
