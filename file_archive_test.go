package dbox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileArchive_simple(t *testing.T) {
	store := NewMemoryStore()

	file := NewFile(store)
	file.SetName("archivefiletest")
	file.SetMapDataStore(store)
	file.SetRawDataStore(store)

	mapSet(file.Meta(), "v1", "1")
	mapSet(file.MapData(), "v2", "2")
	file.RawData().Write([]byte("3"))

	err := file.Sync()
	assert.NoError(t, err, "create file")

	network, err := file.Export()
	assert.NoError(t, err, "export file")

	file = NewFile(store)
	err = file.Import(network)
	assert.NoError(t, err, "import file")

	assert.Equal(t, mapString(file.Meta(), "v1"), "1", "value data file")
	assert.Equal(t, mapString(file.MapData(), "v2"), "2", "value data file")
	assert.Equal(t, file.RawData().Bytes(), []byte("3"), "value data file")
}

func BenchmarkTestFileArchive_simple(b *testing.B) {
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
