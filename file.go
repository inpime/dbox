package dbox

import (
	"github.com/gebv/typed"
    "time"
)

var _ Object = (*File)(nil)

var (
	MapDataIDMetaKey = "MapDataID"
	RawDataIDMetaKey = "RawDataID"
	CreatedAtKey        = "CreatedAt"
	UpdatedAtKey = "UpdatedAt"
    NameKey = "NameKey"
)

func NewFile(store Store) *File {
    
	return &File{
		store:     store,
        MapObject: NewMapObject(store),
	}
}

type File struct {
	*MapObject

	mdata *MapObject
	rdata Object

	store Store

    err error

    // storeRawData Store
    // storeStructData Store
    // storeMetaFileData Store
}

func (f File) String() string {
    return f.Name()
}

func (f File) Error() error {
    return f.err
}

func (f File) Name() string {
    return f.Meta().String(NameKey)
}

func (f File) SetName(v string) {
    f.Meta().Set(NameKey, v)
}

func (f File) UpdatedAt() time.Time {
    return f.Meta().Time(UpdatedAtKey)
}

func (f File) CreatedAt() time.Time {
    return f.Meta().Time(CreatedAtKey)
}

// Meta meta data file
func (f File) Meta() *typed.Typed {
    
	return f.Map()
}

// MapData struct data file
func (f *File) MapData() *typed.Typed {
	if f.mdata == nil {
		f.mdata = NewMapObject(f.store)

		err := f.store.Get(f.mapDataID(), f.mdata)

		if err == ErrNotFound || len(f.mapDataID()) == 0 {
			f.mdata.Sync()
			f.setMapDataID(f.mdata.ID())
			f.MapObject.Sync() // update file props
		} else if err != nil {
			// handler error
            f.err = err

            // TODO: How to address the error?
		}
	}

	return f.mdata.Map()
}

// RawData raw data file
func (f *File) RawData() Object {
	if f.rdata == nil {
		f.rdata = NewRawObject(f.store)

		err := f.store.Get(f.rawDataID(), f.rdata)

		if err == ErrNotFound || len(f.rawDataID()) == 0 {
			f.rdata.Sync()
			f.setRawDataID(f.rdata.ID())
			f.MapObject.Sync() // update file props
		} else if err != nil {
			// handler error
            f.err = err

            // TODO: How to address the error?
		}
	}

	return f.rdata
}

func (f *File) Delete() error {
    if f.IsNew() {
        return ErrEmptyID
    }

    if len(f.Name()) == 0 {
        return ErrEmptyName
    }

    f.MapData()
    if err := f.store.Delete(f.mdata); err != nil {
        return err
    }

    f.RawData()
    if err := f.store.Delete(f.rdata); err != nil {
        return err
    }

    return f.store.Delete(f)
}

func (f *File) Sync() error {
    if f.IsNew() {
        f.id = NewUUID()
        f.BeforeCreate()
    }

    if len(f.Name()) == 0 {
        f.SetName(f.ID())
    }

    f.BeforeUpdate()

	if f.mdata != nil {
		return f.mdata.Sync()
	}

	if f.rdata != nil {
		return f.rdata.Sync()
	}

	// after save related objects

	if err := f.Encode(); err != nil {
		return err
	}

	return f.store.Save(f)
}

//----------
// helpfull
//---------

func (f *File) BeforeCreate() {
	f.Meta().Set(CreatedAtKey, time.Now()) 
}

func (f *File) BeforeUpdate()  {
	f.Meta().Set(UpdatedAtKey, time.Now()) 
}

func (f File) mapDataID() string {
	return f.Meta().String(MapDataIDMetaKey)
}

func (f File) setMapDataID(id string) {
	f.Meta().Set(MapDataIDMetaKey, id)
}

func (f File) rawDataID() string {
	return f.Meta().String(RawDataIDMetaKey)
}

func (f File) setRawDataID(id string) {
	f.Meta().Set(RawDataIDMetaKey, id)
}

// --------

func NewFileID(id string, store Store) (*File, error) {
    file := NewFile(store)
    
	return file, store.Get(id, file) 
}

func NewFileName(name string, store Store) (*File, error) {
    file := NewFile(store)
    
	return file, store.(FileStore).GetByName(name, file)  
}