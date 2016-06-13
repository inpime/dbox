package dbox

import (
	"github.com/boltdb/bolt"
	"path/filepath"
	"strings"
	"time"
)

func NewBoltDBStore(path string) *BoltDBStore {
	store := &BoltDBStore{
		path: path,
	}

	if err := ensureDir(filepath.Dir(path)); err != nil {
		panic(err)
	}

	var err error

	store.db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		panic(err)
	}

	err = store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(store.defaultBucket()))
		return err
	})

	if err != nil {
		panic(err)
	}

	return store
}

type BoltDBStore struct {
	path string
	db   *bolt.DB
}

func (s BoltDBStore) defaultBucket() string {
	return "__dbox.buckets.default"
}

func (s BoltDBStore) Get(id string, obj Object) error {
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.defaultBucket()))
		b := bucket.Get([]byte(strings.ToLower(id)))

		// if empty then it is not found
		// because at the end of must 1 byte (see superDecodeData, superEncodeData)

		if len(b) == 0 {
			return ErrNotFound
		}

		obj.Write(superDecodeData(b))
		obj.SetID(id)
		return nil
	})

	if err != nil {
		return err
	}

	return obj.Decode()
}

func (s *BoltDBStore) save(obj Object) (err error) {
	return s.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.defaultBucket()))
		b := superEncodeData(obj.Bytes())
		return bucket.Put([]byte(strings.ToLower(obj.ID())), b)
	})
}

func (s *BoltDBStore) saveFileRefs(file *File) error {
	// Ref by file name
	ref := NewRefObject(s)
	ref.SetID(file.Name())
	ref.SetRefID(file.ID())

	return s.save(ref)
}

func (s *BoltDBStore) Save(obj Object) (err error) {
	if err := s.save(obj); err != nil {
		return err
	}

	switch obj := obj.(type) {
	case *File:
		err = s.saveFileRefs(obj)
	}

	return
}

func (s *BoltDBStore) delete(obj Object) (err error) {
	return s.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.defaultBucket()))
		return bucket.Delete([]byte(strings.ToLower(obj.ID())))
	})
}

func (s *BoltDBStore) Delete(obj Object) (err error) {
	if err := s.delete(obj); err != nil {
		return err
	}

	switch obj := obj.(type) {
	case *File:
		// load refs file
		fileRef := NewRefObject(s)

		if err = s.Get(obj.Name(), fileRef); err != nil {
			return
		}

		err = s.delete(fileRef)
	}

	return
}

func (s BoltDBStore) Type() StoreType {
	return BoltDBStoreType
}

// GetByName Get object by name (via ReObject)
func (s BoltDBStore) GetByName(name string, obj Object) error {
	ref := NewRefObject(&s)

	if err := s.Get(name, ref); err != nil {
		return err
	}

	return s.Get(ref.RefID(), obj)
}
