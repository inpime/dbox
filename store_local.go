package dbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func NewLocalStore(path string) *LocalStore {
	store := &LocalStore{
		storepath: path,
	}

	if err := ensureDir(path); err != nil {
		panic(err)
	}

	return store
}

type LocalStore struct {
	storepath string
}

func (s LocalStore) formatPathFile(obj Object) string {
	if obj, isFile := obj.(*File); isFile && len(obj.Bucket()) > 0 {
		bucketname := strings.ToLower(obj.Bucket())
		return filepath.Clean(s.storepath + string(os.PathSeparator) + bucketname + string(os.PathSeparator) + obj.ID())
	}

	return filepath.Clean(s.storepath + string(os.PathSeparator) + obj.ID())
}

func (s LocalStore) Get(id string, obj Object) error {
	obj.SetID(id) // to correctly build the path to the file
	filePath := s.formatPathFile(obj)

	if !exists(filePath) {
		obj.SetID("") // force clean, not to be deceived (see function Object.IsNew)
		return ErrNotFound
	}

	b, err := ioutil.ReadFile(filePath)

	// select {
	//     case strings.HasSuffix(err.Error(), "is a directory"):
	//         return ErrInvalidData
	//     default:
	//         return err
	// }

	if err != nil {
		return err
	}

	obj.Write(b)
	obj.SetID(id)

	return obj.Decode()
}

func (s *LocalStore) save(obj Object) error {
	filePath := s.formatPathFile(obj)

	// TODO: если это файл, добавить вконце еще каталог
	// ensureDir(filepath.Dir(filePath))

	return ioutil.WriteFile(filePath, obj.Bytes(), 0644)
}

func (s *LocalStore) Save(obj Object) (err error) {
	if err := s.save(obj); err != nil {
		return err
	}

	switch obj := obj.(type) {
	case *File:
		err = s.saveFileRefs(obj)
	}

	return err
}

func (s *LocalStore) delete(obj Object) error {
	filePath := s.formatPathFile(obj)

	return os.Remove(filePath)
}

func (s *LocalStore) Delete(obj Object) (err error) {
	if err := s.delete(obj); err != nil {
		return err
	}

	switch obj := obj.(type) {
	case *File:
		// load refs file
		fileRef := NewFile(obj.store)
		if err = s.Get(obj.Name(), fileRef); err != nil {
			return
		}

		err = s.delete(fileRef)
	}

	return err
}

func (s LocalStore) Type() StoreType {
	return LocalStoreType
}

// FileStore interface

func (s LocalStore) GetByName(name string, obj Object) error {
	ref := NewRefObject(&s)

	if err := s.Get(name, ref); err != nil {
		return err
	}

	return s.Get(ref.RefID(), obj)
}

func (s *LocalStore) saveFileRefs(file *File) error {
	// Ref by file name
	ref := NewRefObject(s)
	ref.SetID(file.Name())
	ref.SetRefID(file.ID())

	return s.save(ref)
}

//

func exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func ensureDir(filename string) error {
	fdir := path.Dir(filename)
	if fdir != "" && fdir != filename {
		d, err := os.Stat(fdir)
		if err == nil {
			if !d.IsDir() {
				return fmt.Errorf("filename's dir exists but isn't' a directory: filename:%v dir:%v", filename, fdir)
			}
		} else if os.IsNotExist(err) {
			err := os.MkdirAll(fdir, 0775)
			if err != nil {
				return fmt.Errorf("unable to create path. : filename:%v dir:%v err:%v", filename, fdir, err)
			}
		}
	}
	return nil
}
