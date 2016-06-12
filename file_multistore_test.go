package dbox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultiStore_simpleStrategy(t *testing.T) {
	storeFiles := NewMemoryStoreName("a")
	storeMapData := NewMemoryStoreName("b")
	storeRawData := NewMemoryStoreName("c")

	createSimpleStrategy(t, storeFiles, storeMapData, storeRawData)

	assert.Equal(t, len((*MemoryStore)(storeFiles).list), 2, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeMapData).list), 1, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeRawData).list), 1, "count files storage")

	deleteSimpleStrategy(t, storeFiles, storeMapData, storeRawData)

	assert.Equal(t, len((*MemoryStore)(storeFiles).list), 0, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeMapData).list), 0, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeRawData).list), 0, "count files storage")
}

func TestMultiStore_fileAndStorage(t *testing.T) {
	storeFiles := NewMemoryStoreName("a")
	storeMapData := NewMemoryStoreName("b")
	storeRawData := NewMemoryStoreName("c")

	file := NewFile(storeFiles)
	file.SetName("namefile")
	file.SetMapDataStore(storeMapData)
	file.SetRawDataStore(storeRawData)

	file.Map().Set("a", "b")
	file.MapData().Set("c", "d")
	file.RawData().Write([]byte("abc"))

	file.Sync()

	// check store

	assert.Equal(t, len((*MemoryStore)(storeFiles).list), 2, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeMapData).list), 1, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeRawData).list), 1, "count files storage")

	// ----------------
	// load by id
	// ----------------

	file, err := NewFileID(file.ID(), storeFiles)
	assert.NoError(t, err, "get by id")
	file.SetName("namefile")
	file.SetMapDataStore(storeMapData)
	file.SetRawDataStore(storeRawData)

	// before init subobjects
	assert.Equal(t, len((*MemoryStore)(storeFiles).list), 2, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeMapData).list), 1, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeRawData).list), 1, "count files storage")

	assert.Equal(t, file.Meta().String("a"), "b", "not expected value")
	assert.Equal(t, file.MapData().String("c"), "d", "not expected value")
	assert.Equal(t, file.RawData().Bytes(), []byte("abc"), "not expected value")

	// after init subobjects
	assert.Equal(t, len((*MemoryStore)(storeFiles).list), 2, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeMapData).list), 1, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeRawData).list), 1, "count files storage")

	// ----------------
	// load by name
	// ----------------

	file, err = NewFileName("namefile", storeFiles)
	assert.NoError(t, err, "get by name")
	file.SetMapDataStore(storeMapData)
	file.SetRawDataStore(storeRawData)

	// before init subobjects
	assert.Equal(t, len((*MemoryStore)(storeFiles).list), 2, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeMapData).list), 1, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeRawData).list), 1, "count files storage")

	assert.Equal(t, file.Meta().String("a"), "b", "not expected value")
	assert.Equal(t, file.MapData().String("c"), "d", "not expected value")
	assert.Equal(t, file.RawData().Bytes(), []byte("abc"), "not expected value")

	// after init subobjects
	assert.Equal(t, len((*MemoryStore)(storeFiles).list), 2, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeMapData).list), 1, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeRawData).list), 1, "count files storage")

	// ---------------
	// delete file
	// ---------------

	err = file.Delete()
	assert.NoError(t, err, "remove file")

	assert.Equal(t, len((*MemoryStore)(storeFiles).list), 0, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeMapData).list), 0, "count files storage")
	assert.Equal(t, len((*MemoryStore)(storeRawData).list), 0, "count files storage")
}
