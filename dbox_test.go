package dbox

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDbox_simpleStrategy(t *testing.T) {
	InitDefaultStores()
	InitDefaultBuckets()

	file, err := LoadOrNewFile("static", "image.png")
	assert.Equal(t, err, ErrNotFound, "new file")
	assert.Equal(t, file.IsNew(), true, "is new file")
	assert.Equal(t, file.Name(), "image.png", "value file")
	assert.Equal(t, file.Bucket(), "static", "value file")

	file.Meta().Set("ContentType", "image/png")
	file.MapData().Set("a", "b")
	file.RawData().Write([]byte("image data ..."))
	err = file.Sync()
	assert.NoError(t, err, "saved file")

	// Load saved file and check values
	file, err = LoadOrNewFile("static", "image.png")
	assert.NoError(t, err, "load existing file")

	assert.Equal(t, file.Name(), "image.png", "value file")
	assert.Equal(t, file.Bucket(), "static", "value file")
	assert.Equal(t, file.Meta().String("ContentType"), "image/png", "value file")
	assert.Equal(t, file.MapData().String("a"), "b", "value file")
	assert.Equal(t, file.RawData().Bytes(), []byte("image data ..."), "value file")
	assert.Equal(t, file.mdata.store.Type(), MemoryStoreType, "check store type")
	assert.Equal(t, file.rdata.store.Type(), LocalStoreType, "check store type")
	assert.Equal(t, file.store.Type(), MemoryStoreType, "check store type")

	os.RemoveAll(LocalStoreDefaultStorePath)
}
