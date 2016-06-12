package dbox

import (
	"fmt"
)

var (
	ErrNotFound       = fmt.Errorf("not_found")
	ErrNotFoundBucket = fmt.Errorf("not_found_bucket")
	ErrInvalidData    = fmt.Errorf("invalid_data")

	ErrEmptyName = fmt.Errorf("empty_name")
	ErrEmptyID   = fmt.Errorf("empty_id")
)

// BucketStore a store of information about the buckets
var BucketStore Store
var LocalStoreDefaultStorePath = "./workspaces.dbox/"

func InitDefaultStores() {
	RegistryStore("localfs", NewLocalStore(LocalStoreDefaultStorePath))
	RegistryStore("memory", NewMemoryStore())
}

func InitDefaultBuckets() {
	BucketStore = NewMemoryStore()

	bucket := NewBucket()
	bucket.SetName("static")
	bucket.SetRawDataStoreName("localfs")
	bucket.SetMapDataStoreName("memory")
	bucket.SetMetaDataStoreName("memory")

	bucket.Sync()

	bucket = NewBucket()
	bucket.SetName("pages")
	bucket.SetRawDataStoreName("memory")
	bucket.SetMapDataStoreName("memory")
	bucket.SetMetaDataStoreName("memory")

	bucket.Sync()
}
