package dbox

type StoreType string

const (
	MemoryStoreType StoreType = "memory"
	LocalStoreType  StoreType = "local"
	BoltDBStoreType StoreType = "boltdb"
)

type FileStore interface {
	Store

	GetByName(name string, obj Object) error
}

type Store interface {
	Get(id string, obj Object) error
	Save(Object) error
	Delete(Object) error

	Type() StoreType
}
