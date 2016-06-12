package dbox

type StoreType string

const (
	MemoryStoreType StoreType = "memory"
	LocalStoreType  StoreType = "local"
)

// FileStore implements store data of files
type FileStore interface {
	Store

	GetByName(name string, obj Object) error
}

// Store implements store data of objects
type Store interface {
	Get(id string, obj Object) error
	Save(Object) error
	Delete(Object) error

	Type() StoreType
}
