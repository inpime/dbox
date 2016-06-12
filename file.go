package dbox

import (
	"github.com/gebv/typed"
	"time"
)

var _ Object = (*File)(nil)
var _ Object = (*Bucket)(nil)

var (
	NameKey          = "NameKey"
	MapDataIDMetaKey = "MapDataID"
	RawDataIDMetaKey = "RawDataID"
	CreatedAtKey     = "CreatedAt"
	UpdatedAtKey     = "UpdatedAt"

	BucketKey                = "BucketKey"
	MapDataStoreNameKey      = "MapDataStoreNameKey"
	RawDataStoreNameKey      = "RawDataStoreNameKey"
	MetaDataFileStoreNameKey = "MetaDataFileStoreNameKey"
)

type EntityType string

var (
	FileEntityType   EntityType = "file"
	BucketEntityType EntityType = "bucket"
)

func NewBucket() *Bucket {

	return &Bucket{
		File: File{
			store:     BucketStore,
			MapObject: NewMapObject(BucketStore),
		},
	}
}

type Bucket struct {
	File
}

func (Bucket) Type() EntityType {
	return BucketEntityType
}

func (b Bucket) SetMapDataStoreName(v string) {
	b.Meta().Set(MapDataStoreNameKey, v)
}

func (b Bucket) SetRawDataStoreName(v string) {
	b.Meta().Set(RawDataStoreNameKey, v)
}

func (b Bucket) SetMetaDataStoreName(v string) {
	b.Meta().Set(MetaDataFileStoreNameKey, v)
}

func (b Bucket) MapDataStoreName() string {
	return b.Meta().String(MapDataStoreNameKey)
}

func (b Bucket) RawDataStoreName() string {
	return b.Meta().String(RawDataStoreNameKey)
}

func (b Bucket) MetaDataStoreName() string {
	return b.Meta().String(MetaDataFileStoreNameKey)
}

func NewFile(store Store) *File {

	return &File{
		store:     store,
		MapObject: NewMapObject(store),
	}
}

type File struct {
	*MapObject

	mdata *MapObject
	rdata *RawObject

	store Store

	mapDataStore Store
	rawDataStore Store

	invalid       bool
	reasoninvalid error
	isnew         bool
}

func (f *File) SetMapDataStore(s Store) {
	f.mapDataStore = s
}

func (f File) mdataStore() Store {
	if f.mapDataStore != nil {
		return f.mapDataStore
	}

	return f.store
}

func (f *File) mdataObj() *MapObject {
	if f.mdata == nil {
		var err error

		f.mdata = NewMapObject(f.mdataStore())

		if len(f.mapDataID()) != 0 {
			err = f.mdataStore().Get(f.mapDataID(), f.mdata)
		}

		if err == ErrNotFound || len(f.mapDataID()) == 0 {
			f.mdata.Sync()
			f.setMapDataID(f.mdata.ID())
			f.syncOnlyMeta() // update file props
		} else if err != nil {
			// handler error
			f.invalid = true
			f.reasoninvalid = err

			// TODO: How to address the error?
		}
	}

	return f.mdata
}

func (f *File) SetRawDataStore(s Store) {
	f.rawDataStore = s
}

func (f File) rdataStore() Store {
	if f.rawDataStore != nil {
		return f.rawDataStore
	}

	return f.store
}

func (f *File) rdataObj() Object {
	if f.rdata == nil {
		var err error

		f.rdata = NewRawObject(f.rdataStore())

		if len(f.rawDataID()) != 0 {
			err = f.rdataStore().Get(f.rawDataID(), f.rdata)
		}

		if err == ErrNotFound || len(f.rawDataID()) == 0 {
			f.rdata.Sync()
			f.setRawDataID(f.rdata.ID())
			f.syncOnlyMeta() // update file props
		} else if err != nil {
			// handler error
			f.invalid = true
			f.reasoninvalid = err

			// TODO: How to address the error?
		}
	}

	return f.rdata
}

func (f File) String() string {
	return f.Name()
}

func (File) Type() EntityType {
	return FileEntityType
}

func (f File) Bucket() string {
	return f.Meta().String(BucketKey)
}

func (f File) SetBucket(v string) {
	f.Meta().Set(BucketKey, v)
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
func (f *File) Meta() *typed.Typed {

	return f.MapObject.Map()
}

// MapData struct data file
func (f *File) MapData() *typed.Typed {

	return f.mdataObj().Map()
}

// RawData raw data file
func (f *File) RawData() Object {

	return f.rdataObj()
}

func (f *File) Delete() error {
	if f.invalid {
		return f.reasoninvalid
	}

	if f.IsNew() {
		return ErrEmptyID
	}

	if len(f.Name()) == 0 {
		return ErrEmptyName
	}

	if err := f.mdataStore().Delete(f.mdataObj()); err != nil {
		return err
	}

	if err := f.rdataStore().Delete(f.rdataObj()); err != nil {
		return err
	}

	return f.store.Delete(f)
}

func (f *File) syncOnlyMeta() error {
	if f.IsNew() {
		f.id = NewUUID()
		f.BeforeCreate()
	}

	f.BeforeUpdate()

	return f.MapObject.Sync()
}

func (f *File) Sync() error {
	if f.invalid {
		return f.reasoninvalid
	}

	if err := f.syncOnlyMeta(); err != nil {
		return err
	}

	//

	if f.mdata != nil {
		if err := f.mdataObj().Sync(); err != nil {
			return err
		}
	}

	if f.rdata != nil {
		if err := f.rdataObj().Sync(); err != nil {
			return err
		}
	}

	return f.store.Save(f)
}

//----------
// helpfull
//---------

func (f *File) BeforeCreate() {
	f.Meta().Set(CreatedAtKey, time.Now())
}

func (f *File) BeforeUpdate() {
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

func NewFileID(id string, store Store) (file *File, err error) {
	file = NewFile(store)

	err = store.Get(id, file)

	if err == nil {
		// file.init()
	}

	return
}

func NewFileName(name string, store Store) (file *File, err error) {
	file = NewFile(store)

	err = store.(FileStore).GetByName(name, file)

	if err == nil {
		// file.init()
	}

	return
}

func LoadOrNewFile(bucketName string, fileName string) (file *File, err error) {
	bucket, err := BucketByName(bucketName)
	if err != nil {

		if err == ErrNotFound {
			return nil, ErrNotFoundBucket
		}

		return
	}

	file, err = NewFileName(fileName, MustStore(bucket.MetaDataStoreName()))
	file.SetMapDataStore(MustStore(bucket.MapDataStoreName()))
	file.SetRawDataStore(MustStore(bucket.RawDataStoreName()))

	if err == ErrNotFound {
		file.SetName(fileName)
		file.SetBucket(bucketName)
	}

	return
}

func BucketByName(name string) (file *Bucket, err error) {
	file = NewBucket()
	err = BucketStore.(FileStore).GetByName(name, file)

	return
}
