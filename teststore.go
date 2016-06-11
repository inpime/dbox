package dbox

import (
	"sync"
)


func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		list: make(map[string][]byte),
	}
}

type MemoryStore struct {
	sync.RWMutex

	list map[string][]byte
}

// Store interface

func (s MemoryStore) Get(id string, obj Object) error {
	s.Lock()
	defer s.Unlock()

	if b, exists := s.list[id]; exists {
		obj.SetID(id)
		obj.Write(b)

		return obj.Decode()
	}

	return ErrNotFound
}

func (s *MemoryStore) save(obj Object) error {
    s.Lock()
	defer s.Unlock()

	s.list[obj.ID()] = obj.Bytes()

    return nil
}

func (s *MemoryStore) Save(obj Object) (err error) {
    
    if err := s.save(obj); err != nil {
        return err
    }

    switch obj := obj.(type) {
        case *File:
            err = s.saveFileRefs(obj)
    }

	return err
}

func (s *MemoryStore) delete(id string) error {
    s.Lock()
	defer s.Unlock()

	delete(s.list, id)

    return nil
}

func (s *MemoryStore) Delete(obj Object) (err error) {
    if err := s.delete(obj.ID()); err != nil {
        return err
    }

    switch obj := obj.(type) {
        case *File:
            err = s.delete(obj.Name())
    }

    return err
}

func (s MemoryStore) Type() StoreType {
	return MemoryStoreType
}

// FileStore interface

func (s MemoryStore) GetByName(name string, obj Object) error {
    ref := NewRefObject(&s)

    if err := s.Get(name, ref); err != nil {
        return err
    }

    return s.Get(ref.RefID(), obj) 
}

func (s *MemoryStore) saveFileRefs(file *File) error {
    // Ref by file name
    ref := NewRefObject(s)
    ref.SetID(file.Name())
    ref.SetRefID(file.ID())
    
    return s.save(ref)
}